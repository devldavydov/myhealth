import streamlit as st
import requests

DATA_KEY = "food.edit.data"

def handler():
    key = st.query_params.get("key", "")

    # Load from server
    if DATA_KEY in st.session_state:
        food = st.session_state[DATA_KEY]
    else:    
        try:
            resp = requests.get(f"{st.session_state["api_url"]}/food/{key}").json()
        except Exception as ex:
            st.markdown(f"### Редактирование еды")
            st.error(ex)
            return

        if resp["error"] != "":
            st.markdown(f"### Редактирование еды")
            st.error(resp["error"])
            return

        food = resp["data"]
        st.session_state[DATA_KEY] = food

    # Create page
    st.markdown(f'### Редактирование еды "{food["name"]}"')
   
    name = st.text_input("Наименование", value=food["name"])
    brand = st.text_input("Бренд", value=food["brand"])
    cal100 = st.number_input(
        "ККал, 100г.",
        value=float(food["cal100"]),
        min_value=0.0,
        format="%0.2f",
        step=0.01)
    prot100 = st.number_input(
        "Белки, 100г.",
        value=float(food["prot100"]),
        min_value=0.0,
        format="%0.2f",
        step=0.01)
    fat100 = st.number_input(
        "Жиры, 100г.",
        value=float(food["fat100"]),
        min_value=0.0,
        format="%0.2f",
        step=0.01)
    carb100 = st.number_input(
        "Углеводы, 100г.",
        value=float(food["carb100"]),
        min_value=0.0,
        format="%0.2f",
        step=0.01,)
    comment = st.text_area("Комментарий", value=food["comment"])

    if st.button("Сохранить", icon=":material/check:", type="primary"):
        req = {
            "key": st.session_state[DATA_KEY]["key"],
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
            st.error(ex)
            return

        if resp["error"] != "":
            st.error(resp["error"])
            return

        del(st.session_state[DATA_KEY])
        st.success("Изменения сохранены")

handler()