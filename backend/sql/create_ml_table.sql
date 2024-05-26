DECLARE discount_rate FLOAT64 DEFAULT 0.15;

CREATE TEMP FUNCTION calculate_fcf_cagr( start_fcf FLOAT64, end_fcf FLOAT64, cagr_years INT64 ) RETURNS FLOAT64 AS ( CASE WHEN start_fcf > 0 AND end_fcf > 0 THEN ROUND((POW(start_fcf / end_fcf, 1 / cagr_years) - 1), 4) WHEN start_fcf > end_fcf THEN ROUND((POW((start_fcf - end_fcf + ABS(end_fcf)) / ABS(end_fcf), 1 / cagr_years) - 1), 4) WHEN (start_fcf < 0 AND end_fcf < 0) AND (start_fcf < end_fcf) THEN ROUND((POW(ABS(start_fcf) / ABS(end_fcf), 1 / cagr_years) - 1), 4) * -1 WHEN start_fcf < 0 AND end_fcf > 0 THEN ROUND(((start_fcf / end_fcf) / cagr_years), 4) ELSE NULL END );

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

DROP TABLE IF EXISTS `${GOOGLE_CLOUD_PROJECT}.${DATASET_NAME}.${ML_TABLE_NAME}`;
CREATE TABLE IF NOT EXISTS `${GOOGLE_CLOUD_PROJECT}.${DATASET_NAME}.${ML_TABLE_NAME}` ( date date, ticker STRING, fcf FLOAT64, fcf_cagr_5yrs FLOAT64, assets INT64, intrinsic_val FLOAT64, mrkt_cap FLOAT64, mrkt_intrinsic_ratio FLOAT64, stockprice FLOAT64, stock_5yr_growth FLOAT64, volume_lst_yr FLOAT64, stockprice_future_1yr FLOAT64 );
INSERT INTO `${GOOGLE_CLOUD_PROJECT}.${DATASET_NAME}.${ML_TABLE_NAME}` (date, ticker, fcf, fcf_cagr_5yrs, assets, intrinsic_val, mrkt_cap, mrkt_intrinsic_ratio, stockprice, stock_5yr_growth, volume_lst_yr, stockprice_future_1yr)

WITH tfcf AS
(
    SELECT  cy + 1                               AS cy,
            ticker,
            netcash - propertyexp                AS fcf,
            shares,
            cashasset + investments + securities AS assets
    FROM `${DATASET_NAME}.${FINANICIAL_TABLE_NAME}`
    WHERE cy BETWEEN 2008 AND 2023
    AND netcash <> 0
    AND propertyexp > 0
    AND shares > 0
), ca AS
(
    SELECT  a.cy,
            a.ticker,
            a.shares,
            a.fcf,
            calculate_fcf_cagr(a.fcf,b.fcf,5) AS fcf_cagr_5yrs,
            a.assets
    FROM tfcf a
    INNER JOIN tfcf b
    ON a.ticker = b.ticker AND a.cy = b.cy + 5
), sp AS
(
    SELECT  datetime,
            EXTRACT( YEAR FROM datetime ) AS yr,
            ticker,
            CLOSE AS stockprice
    FROM `${DATASET_NAME}.${STOCK_TABLE_NAME}_monthly`
    WHERE datetime BETWEEN "2009-01-01" AND "2024-01-01"
    AND EXTRACT( MONTH FROM datetime ) = 1
), spf AS
(
    SELECT  a.datetime,
            a.yr,
            a.ticker,
            a.stockprice,
            ROUND(((a.stockprice / b.stockprice) - 1),4) AS stock_5yr_growth,
            c.stockprice                                 AS stockprice_future_1yr
    FROM sp AS a
    INNER JOIN sp AS b
    ON a.ticker = b.ticker AND a.yr = b.yr + 5
    LEFT JOIN sp AS c
    ON a.ticker = c.ticker AND a.yr = c.yr - 1
), volumes AS
(
    SELECT  ticker,
            EXTRACT( YEAR FROM datetime ) + 1 AS yr,
            SUM(volume) AS volume
    FROM `${DATASET_NAME}.${STOCK_TABLE_NAME}_monthly`
    WHERE datetime BETWEEN "2009-01-01" AND "2023-12-01"
    GROUP BY  ticker,
              yr
), con AS
(
    SELECT  CAST(s.datetime AS date) AS date,
            s.ticker,
            c.fcf,
            c.fcf_cagr_5yrs,
            c.assets,
            calculate_intrinsic_value(c.fcf,c.fcf_cagr_5yrs,c.assets,discount_rate) AS intrinsic_val,
            s.stockprice * c.shares                                                 AS mrkt_cap,
            s.stockprice,
            s.stock_5yr_growth,
            v.volume AS volume_lst_yr,
            s.stockprice_future_1yr
    FROM spf AS s
    INNER JOIN ca AS c
    ON s.ticker = c.ticker AND s.yr = c.cy
    LEFT JOIN volumes AS v
    ON s.ticker = v.ticker AND s.yr = v.yr
)
SELECT  date,
        ticker,
        fcf,
        fcf_cagr_5yrs,
        assets,
        intrinsic_val,
        mrkt_cap,
        mrkt_cap/intrinsic_val AS mrkt_intrinsic_ratio,
        stockprice,
        stock_5yr_growth,
        volume_lst_yr,
        stockprice_future_1yr
FROM con
