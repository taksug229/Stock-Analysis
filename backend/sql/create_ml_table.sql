DECLARE discount_rate FLOAT64 DEFAULT 0.15;

CREATE TEMP FUNCTION calculate_cagr( start_fcf FLOAT64, end_fcf FLOAT64, cagr_years INT64 ) RETURNS FLOAT64 AS ( CASE WHEN start_fcf = 0 OR end_fcf = 0 THEN NULL WHEN start_fcf > 0 AND end_fcf > 0 THEN ROUND((POW(start_fcf / end_fcf, 1 / cagr_years) - 1), 4) WHEN start_fcf > end_fcf THEN ROUND((POW((start_fcf - end_fcf + ABS(end_fcf)) / ABS(end_fcf), 1 / cagr_years) - 1), 4) WHEN (start_fcf < 0 AND end_fcf < 0) AND (start_fcf < end_fcf) THEN ROUND((POW(ABS(start_fcf) / ABS(end_fcf), 1 / cagr_years) - 1), 4) * -1 WHEN start_fcf < 0 AND end_fcf > 0 THEN ROUND(((start_fcf / end_fcf) / cagr_years), 4) ELSE NULL END );

CREATE TEMP FUNCTION check_valid_cagr( cy ANY TYPE, revenue ANY TYPE, fcf ANY TYPE, assets ANY TYPE ) RETURNS BOOL AS ( cy IS NOT NULL AND COALESCE(revenue, 0) != 0 AND COALESCE(fcf, 0) != 0 AND COALESCE(assets, 0) != 0 );

CREATE TEMP FUNCTION calculate_intrinsic_value( free_cash_flow FLOAT64, growth_rate FLOAT64, assets FLOAT64, discount_rate FLOAT64 ) RETURNS FLOAT64 AS ( (
WITH compounded AS
(
    SELECT  i                                       AS period,
            free_cash_flow * POW(1 + growth_rate,i) AS comp_value
    FROM UNNEST
    (GENERATE_ARRAY(1, 10)
    ) AS i
), discounted AS
(
    SELECT  period,
            comp_value,
            comp_value * (1 / POW(1 + discount_rate,period)) AS discounted_value
    FROM compounded
), terminal_value AS
(
    SELECT  ((
    SELECT  comp_value
    FROM compounded
    WHERE period = 10) * 10) * (1 / POW(1 + discount_rate, 10)) AS discounted_terminal_value
)
SELECT  (ROUND(SUM(discounted_value) + (
SELECT  discounted_terminal_value
FROM terminal_value) + assets, 0))
FROM discounted ) );

DROP TABLE IF EXISTS `${DATASET_NAME}.${ML_TABLE_NAME}`;
CREATE TABLE IF NOT EXISTS `${DATASET_NAME}.${ML_TABLE_NAME}` ( date date, ticker STRING, revenue INT64, fcf INT64, assets INT64, cagr_yrs INT64, revenue_cagr FLOAT64, fcf_cagr FLOAT64, assets_cagr FLOAT64, intrinsic_val BIGNUMERIC, mrkt_cap INT64, mrkt_intrinsic_ratio FLOAT64, stockprice_last_yr FLOAT64, stockprice_current FLOAT64, stock_cagr FLOAT64, volume_last_yr INT64, stockprice_future_1yr FLOAT64 );
INSERT INTO `${DATASET_NAME}.${ML_TABLE_NAME}` (date, ticker, revenue, fcf, assets, cagr_yrs, revenue_cagr, fcf_cagr, assets_cagr, intrinsic_val, mrkt_cap, mrkt_intrinsic_ratio, stockprice_last_yr, stockprice_current, stock_cagr, volume_last_yr, stockprice_future_1yr)
WITH tfcf AS
(
    SELECT  cy + 1                AS cy,
            ticker,
            revenue,
            netcash - propertyexp AS fcf,
            shares,
            CASE WHEN (investments = securities AND investments > 0) THEN CAST(cashasset + investments AS INT64)  ELSE CAST(cashasset + investments + securities AS INT64) END AS assets
    FROM `${DATASET_NAME}.${FINANICIAL_TABLE_NAME}`
    WHERE cy BETWEEN 2008 AND 2023
    AND propertyexp > 0
    AND shares > 0
    AND (cashasset > 0 OR investments > 0 OR securities > 0 )
), ca AS
(
    SELECT  *
    FROM
    (
        SELECT  a.cy,
                a.ticker,
                a.shares,
                a.revenue,
                a.fcf,
                a.cy - CASE WHEN check_valid_cagr(b.cy,b.revenue,b.fcf,b.assets) THEN b.cy WHEN check_valid_cagr(c.cy,c.revenue,c.fcf,c.assets ) THEN c.cy WHEN check_valid_cagr(d.cy,d.revenue,d.fcf,d.assets) THEN d.cy WHEN check_valid_cagr(e.cy,e.revenue,e.fcf,e.assets) THEN e.cy ELSE NULL END AS cagr_yrs,
                CASE WHEN check_valid_cagr(b.cy,b.revenue,b.fcf,b.assets) THEN calculate_cagr(a.revenue,b.revenue,10)
                     WHEN check_valid_cagr(c.cy,c.revenue,c.fcf,c.assets ) THEN calculate_cagr(a.revenue,c.revenue,7)
                     WHEN check_valid_cagr(d.cy,d.revenue,d.fcf,d.assets) THEN calculate_cagr(a.revenue,d.revenue,5)
                     WHEN check_valid_cagr(e.cy,e.revenue,e.fcf,e.assets) THEN calculate_cagr(a.revenue,e.revenue,3) END AS revenue_cagr,
                CASE WHEN check_valid_cagr(b.cy,b.revenue,b.fcf,b.assets) THEN calculate_cagr(a.fcf,b.fcf,10)
                     WHEN check_valid_cagr(c.cy,c.revenue,c.fcf,c.assets ) THEN calculate_cagr(a.fcf,c.fcf,7)
                     WHEN check_valid_cagr(d.cy,d.revenue,d.fcf,d.assets) THEN calculate_cagr(a.fcf,d.fcf,5)
                     WHEN check_valid_cagr(e.cy,e.revenue,e.fcf,e.assets) THEN calculate_cagr(a.fcf,e.fcf,3) END         AS fcf_cagr,
                CASE WHEN check_valid_cagr(b.cy,b.revenue,b.fcf,b.assets) THEN calculate_cagr(a.assets,b.assets,10)
                     WHEN check_valid_cagr(c.cy,c.revenue,c.fcf,c.assets ) THEN calculate_cagr(a.assets,c.assets,7)
                     WHEN check_valid_cagr(d.cy,d.revenue,d.fcf,d.assets) THEN calculate_cagr(a.assets,d.assets,5)
                     WHEN check_valid_cagr(e.cy,e.revenue,e.fcf,e.assets) THEN calculate_cagr(a.assets,e.assets,3) END   AS assets_cagr,
                a.assets
        FROM tfcf a
        LEFT JOIN tfcf b
        ON a.ticker = b.ticker AND a.cy = b.cy + 10
        LEFT JOIN tfcf c
        ON a.ticker = c.ticker AND a.cy = c.cy + 7
        LEFT JOIN tfcf d
        ON a.ticker = d.ticker AND a.cy = d.cy + 5
        LEFT JOIN tfcf e
        ON a.ticker = e.ticker AND a.cy = e.cy + 3
    )
    WHERE cagr_yrs IS NOT NULL
    AND shares > 0
), sp AS
(
    SELECT  datetime,
            EXTRACT( YEAR
    FROM datetime ) AS yr, ticker, open AS stockprice
    FROM `${DATASET_NAME}.${STOCK_TABLE_NAME}_monthly`
    WHERE datetime BETWEEN "2009-01-01" AND "2024-01-01"
    AND EXTRACT( MONTH
    FROM datetime ) = 1
), spf AS
(
    SELECT  *
    FROM
    (
        SELECT  a.datetime,
                a.yr,
                a.ticker,
                a.stockprice                                                                            AS stockprice_current,
                f.stockprice                                                                            AS stockprice_last_yr,
                a.yr - CASE WHEN b.yr IS NOT NULL THEN b.yr WHEN c.yr IS NOT NULL THEN c.yr WHEN d.yr IS NOT NULL THEN d.yr WHEN e.yr IS NOT NULL THEN e.yr END AS cagr_yrs,
                CASE WHEN b.stockprice IS NOT NULL THEN calculate_cagr(a.stockprice,b.stockprice,10)
                     WHEN c.stockprice IS NOT NULL THEN calculate_cagr(a.stockprice,c.stockprice,7)
                     WHEN d.stockprice IS NOT NULL THEN calculate_cagr(a.stockprice,d.stockprice,5)
                     WHEN e.stockprice IS NOT NULL THEN calculate_cagr(a.stockprice,e.stockprice,3) END AS stock_cagr,
                g.stockprice                                                                            AS stockprice_future_1yr
        FROM sp AS a
        LEFT JOIN sp AS b
        ON a.ticker = b.ticker AND a.yr = b.yr + 10
        LEFT JOIN sp AS c
        ON a.ticker = c.ticker AND a.yr = c.yr + 7
        LEFT JOIN sp AS d
        ON a.ticker = d.ticker AND a.yr = d.yr + 5
        LEFT JOIN sp AS e
        ON a.ticker = e.ticker AND a.yr = e.yr + 3
        LEFT JOIN sp AS f
        ON a.ticker = f.ticker AND a.yr = f.yr + 1
        LEFT JOIN sp AS g
        ON a.ticker = g.ticker AND a.yr = g.yr - 1
    )
    WHERE cagr_yrs IS NOT NULL
), volumes AS
(
    SELECT  ticker,
            EXTRACT( YEAR
    FROM datetime ) + 1 AS yr, SUM(volume) AS volume
    FROM `${DATASET_NAME}.${STOCK_TABLE_NAME}_monthly`
    WHERE datetime BETWEEN "2009-01-01" AND "2023-12-01"
    GROUP BY  ticker,
              yr
), con AS
(
    SELECT  CAST(s.datetime AS date)                                                    AS date,
            s.ticker,
            c.revenue,
            c.fcf,
            c.assets,
            c.cagr_yrs,
            c.revenue_cagr,
            c.fcf_cagr,
            c.assets_cagr,
            ROUND(calculate_intrinsic_value(c.fcf,c.fcf_cagr,c.assets,discount_rate),0) AS intrinsic_val,
            CAST(ROUND(s.stockprice_current * c.shares,0) AS INT64)                     AS mrkt_cap,
            s.stockprice_last_yr,
            s.stockprice_current,
            s.stock_cagr,
            CAST(v.volume AS INT64)                                                     AS volume_last_yr,
            s.stockprice_future_1yr
    FROM spf AS s
    INNER JOIN ca AS c
    ON s.ticker = c.ticker AND s.yr = c.cy
    LEFT JOIN volumes AS v
    ON s.ticker = v.ticker AND s.yr = v.yr
)
SELECT  date,
        ticker,
        revenue,
        fcf,
        assets,
        cagr_yrs,
        revenue_cagr,
        fcf_cagr,
        assets_cagr,
        CAST(CASE WHEN intrinsic_val > 0 THEN intrinsic_val ELSE 0 END                        AS BIGNUMERIC)AS intrinsic_val,
        mrkt_cap,
        LEAST(ROUND(CASE WHEN intrinsic_val > 0 THEN mrkt_cap/intrinsic_val ELSE 0 END,4),10) AS mrkt_intrinsic_ratio,
        stockprice_last_yr,
        stockprice_current,
        stock_cagr,
        volume_last_yr,
        stockprice_future_1yr
FROM con
ORDER BY date, ticker
