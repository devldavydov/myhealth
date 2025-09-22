import streamlit as st
import pandas as pd

def handler():
    st.markdown("### Управление едой")
    st.dataframe(
        pd.DataFrame({
            "link": ["/food-edit?key=ммс"],
            "Ключ": ["ммс"],
            "Наименование": ["MMS"],
            "Бренд": ["Mars"],
            "ККал, 100г.": [200],
            "БЖУ, 100г.": ["12/12/12"],
            "Комментарий": [""]
        }),
        hide_index=True,
        column_config={
            "link": st.column_config.LinkColumn(
                "",
                width="small",
                pinned=True,
                display_text=":material/edit:"
            )
        }
    )

handler()