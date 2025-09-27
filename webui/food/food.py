import streamlit as st
import pandas as pd
import requests

def handler():
    st.markdown("### Управление едой")

    # Load from server
    try:
        resp = requests.get(f"{st.session_state["api_url"]}/food").json()
    except Exception as ex:
        st.error(ex)
        return

    if resp["error"] != "":
        st.error(resp["error"])
        return
    
    # Create page
    if len(resp["data"]) == 0:
        st.info("Список еды пуст")
    else:
        link_col, name_col, brand_col = [], [], []
        cal_col, pfc_col, comment_col = [], [], []

        for item in resp["data"]:
            link_col.append(f"/food-edit?key={item["key"]}")
            name_col.append(item["name"])
            brand_col.append(item["brand"])
            cal_col.append(item["cal100"])
            pfc_col.append(f"{item["prot100"]:.2f} / {item["fat100"]:.2f} / {item["carb100"]:.2f}")
            comment_col.append(item["comment"])

        st.dataframe(
            pd.DataFrame({
                "link": link_col,
                "Наименование": name_col,
                "Бренд": brand_col,
                "ККал, 100г.": cal_col,
                "БЖУ, 100г.": pfc_col,
                "Комментарий": comment_col
            }),
            hide_index=True,
            column_config={
                "link": st.column_config.LinkColumn(
                    "",
                    width="small",
                    pinned=True,
                    display_text="Ред."
                )
            }
        )

    if st.button("Добавить", icon=":material/add:", type="primary"):
        st.switch_page("food/create.py")

handler()