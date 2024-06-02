SELECT  intrinsic_val,
        mrkt_cap,
        predicted_stockprice_future_1yr
FROM ML.PREDICT
(MODEL `${DATASET_NAME}.ml_model`, (
    SELECT  a.date,
            a.ticker,
            a.revenue,
            a.fcf,
            a.assets,
            a.cagr_yrs,
            a.revenue_cagr,
            a.fcf_cagr,
            a.assets_cagr,
            a.intrinsic_val,
            CAST(ROUND(b.shares * ${LIVE_STOCK_PRICE},0) AS INT64)                                                            AS mrkt_cap,
            LEAST(ROUND(CASE WHEN a.intrinsic_val > 0 THEN (b.shares * ${LIVE_STOCK_PRICE})/a.intrinsic_val ELSE 0 END,4),10) AS mrkt_intrinsic_ratio,
            a.stockprice_last_yr,
            ${LIVE_STOCK_PRICE} AS stockprice_current,
            a.stock_cagr,
            a.volume_last_yr
    FROM
    (
        SELECT  *
        FROM `${DATASET_NAME}.${ML_TABLE_NAME}`
        WHERE ticker = "${TICKER}"
        AND date = "2024-01-01"
    ) AS a
    LEFT JOIN
    (
        SELECT  ticker,
                shares,

        FROM `${DATASET_NAME}.${FINANICIAL_TABLE_NAME}`
        WHERE ticker = "${TICKER}"
        AND cy = 2023
    ) AS b
    ON a.ticker = b.ticker )
)
