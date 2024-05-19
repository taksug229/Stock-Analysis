import pandas as pd
import numpy as np

from sklearn.tree import DecisionTreeRegressor
from sklearn.metrics import mean_squared_error

file_path = "data/stock_price_monthly.csv"


def main():
    df = pd.read_csv(file_path)
    df["datetime"] = pd.to_datetime(df["datetime"])
    TARGET = "close"
    del df["ticker"]
    X_train = df.query("datetime <= '2023-12-31'").set_index("datetime")
    X_test = df.query("datetime > '2023-12-31'").set_index("datetime")
    y_train = X_train.pop(TARGET)
    y_test = X_test.pop(TARGET)
    model_ = DecisionTreeRegressor()
    model = model_.fit(X_train, y_train)
    y_pred = model.predict(X_test)
    rmse = np.sqrt(mean_squared_error(y_test, y_pred))
    mean_close = np.mean(y_test)
    print(f"Mean Close: {mean_close}")
    print(f"RMSE: {rmse} (Error rate: {rmse/mean_close})")


if __name__ == "__main__":
    main()
