SELECT  a.intrinsic_val,
        b.shares
FROM
(
    SELECT  ticker,
            intrinsic_val
    FROM `${DATASET_NAME}.${ML_TABLE_NAME}`
    WHERE ticker = "${TICKER}"
    AND date = "2024-01-01"
) AS a
LEFT JOIN
(
    SELECT  ticker,
            shares
    FROM `${DATASET_NAME}.${FINANICIAL_TABLE_NAME}`
    WHERE ticker = "${TICKER}"
    AND cy = 2023
) AS b
ON a.ticker = b.ticker
