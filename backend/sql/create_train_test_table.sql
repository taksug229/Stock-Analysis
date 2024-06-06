DROP TABLE IF EXISTS `${DATASET_NAME}.train_data`;
CREATE TABLE IF NOT EXISTS `${DATASET_NAME}.train_data` ( date date, ticker STRING, revenue INT64, fcf INT64, assets INT64, cagr_yrs INT64, revenue_cagr FLOAT64, fcf_cagr FLOAT64, assets_cagr FLOAT64, intrinsic_val BIGNUMERIC, mrkt_cap INT64, mrkt_intrinsic_ratio FLOAT64, stockprice_last_yr FLOAT64, stockprice_current FLOAT64, stock_cagr FLOAT64, 52w_high FLOAT64, 52w_low FLOAT64, volume_last_yr INT64, stockprice_future_1yr FLOAT64 );
INSERT INTO `${DATASET_NAME}.train_data` (date, ticker, revenue, fcf, assets, cagr_yrs, revenue_cagr, fcf_cagr, assets_cagr, intrinsic_val, mrkt_cap, mrkt_intrinsic_ratio, stockprice_last_yr, stockprice_current, stock_cagr, 52w_high, 52w_low, volume_last_yr, stockprice_future_1yr)
SELECT  *
FROM stocks.ml
WHERE date <= "2022-01-01"
ORDER BY date, ticker;

DROP TABLE IF EXISTS `${DATASET_NAME}.test_data`;
CREATE TABLE IF NOT EXISTS `${DATASET_NAME}.test_data` ( date date, ticker STRING, revenue INT64, fcf INT64, assets INT64, cagr_yrs INT64, revenue_cagr FLOAT64, fcf_cagr FLOAT64, assets_cagr FLOAT64, intrinsic_val BIGNUMERIC, mrkt_cap INT64, mrkt_intrinsic_ratio FLOAT64, stockprice_last_yr FLOAT64, stockprice_current FLOAT64, stock_cagr FLOAT64, 52w_high FLOAT64, 52w_low FLOAT64, volume_last_yr INT64, stockprice_future_1yr FLOAT64 );
INSERT INTO `${DATASET_NAME}.test_data` (date, ticker, revenue, fcf, assets, cagr_yrs, revenue_cagr, fcf_cagr, assets_cagr, intrinsic_val, mrkt_cap, mrkt_intrinsic_ratio, stockprice_last_yr, stockprice_current, stock_cagr, 52w_high, 52w_low, volume_last_yr, stockprice_future_1yr)
SELECT  *
FROM stocks.ml
WHERE date = "2023-01-01"
ORDER BY date, ticker;
