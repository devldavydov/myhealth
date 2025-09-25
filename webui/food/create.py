import streamlit as st
import requests
import uuid

def handler():
    st.markdown(f"### Создание еды")
   
    name = st.text_input("Наименование")
    brand = st.text_input("Бренд")
    cal100 = st.number_input(
        "ККал, 100г.",
        min_value=0.0,
        format="%0.2f",
        step=0.01)
    prot100 = st.number_input(
        "Белки, 100г.",
        min_value=0.0,
        format="%0.2f",
        step=0.01)
    fat100 = st.number_input(
        "Жиры, 100г.",
        min_value=0.0,
        format="%0.2f",
        step=0.01)
    carb100 = st.number_input(
        "Углеводы, 100г.",
        min_value=0.0,
        format="%0.2f",
        step=0.01,)
    comment = st.text_area("Комментарий")

    if st.button("Сохранить", icon=":material/check:", type="primary"):
        req = {
            "key": str(uuid.uuid4()),
            "name": name,
            "brand": brand,
            "cal100": cal100,
            "prot100": prot100,
            "fat100": fat100,
            "carb100": carb100,
            "comment": comment
        }

        try:
            resp = requests.post(f"{st.session_state["api_url"]}/food/set", json=req).json()
        except Exception as ex:
            st.toast(f":red[{ex}]", duration="long", icon=":material/error:")
            return

        if resp["error"] != "":
            st.toast(f":red[{resp["error"]}]", duration="long", icon=":material/error:")
            return

        st.switch_page("food/food.py")

handler()