CREATE OR REPLACE MODEL `${DATASET_NAME}.ml_model` OPTIONS(model_type = 'AUTOML_REGRESSOR', OPTIMIZATION_OBJECTIVE = 'MINIMIZE_RMSE', input_label_cols = ['stockprice_future_1yr'], budget_hours = 1.0) AS
SELECT  *
FROM `${DATASET_NAME}.train_data`
