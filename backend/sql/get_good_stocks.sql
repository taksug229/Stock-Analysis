SELECT DISTINCT ticker, mrkt_cap
FROM `${DATASET_NAME}.${ML_TABLE_NAME}`
WHERE date = "2024-01-01"
AND mrkt_intrinsic_ratio BETWEEN 0 AND 0.7
ORDER BY mrkt_cap DESC
