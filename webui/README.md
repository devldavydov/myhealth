# Install

```
sudo apt install python3-venv
python3 -m venv ./venv
source venv/bin/activate

pip install -r requirements.txt
```

# Run
```
source venv/bin/activate
streamlit run app.py --server.port=9090 --server.address=0.0.0.0
```

# Deactivate venv
```
deactivate
```