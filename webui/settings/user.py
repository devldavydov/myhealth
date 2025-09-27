import streamlit as st
import requests

DATA_KEY = "usersettings.edit.data"

def get_default_user_settings():
    return {"calLimit": 0}

def handler():
    # Load from server
    if DATA_KEY in st.session_state:
        us = st.session_state[DATA_KEY]
    else:    
        try:
            resp = requests.get(f"{st.session_state["api_url"]}/usersettings")
            resp_obj = resp.json()
        except Exception as ex:
            st.markdown(f"### Настройки пользователя")
            st.error(ex)
            return

        if resp_obj["error"] != "" and resp.status_code == 200:
            st.markdown(f"### Настройки пользователя")
            st.error(resp["error"])
            return
        
        if resp_obj["error"] != "" and resp.status_code == 404:
            us = get_default_user_settings()
            st.toast(resp_obj["error"])
        else:
            us = resp_obj["data"]
        
        st.session_state[DATA_KEY] = us

    st.markdown("### Настройки пользователя")

    calLimit = st.number_input(
        "Лимит ККал",
        value=float(us["calLimit"]),
        min_value=0.0,
        format="%0.2f",
        step=0.01)
     
    if st.button("Сохранить", icon=":material/check:", type="primary"):
        pass

handler()