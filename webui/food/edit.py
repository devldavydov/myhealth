import streamlit as st

def handler():
    key = st.query_params["key"]
    st.markdown("### Редактирование еды " + key)

handler()