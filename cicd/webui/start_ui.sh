#!/bin/bash

source venv/bin/activate
streamlit run app.py --server.port=9090 --server.address=192.168.100.100