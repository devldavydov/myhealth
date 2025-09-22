import streamlit as st
import pandas as pd

st.markdown("### Расчет лимита ккал")

gender = st.radio(
    "**Пол**",
    [":material/man: Мужчина", ":material/woman: Женщина"])

weight = st.number_input(
    "**Вес, кг.**",
    min_value=0.0,
    step=0.1,
    format="%0.1f")

height = st.number_input(
    "**Рост, см.**",
    min_value=0.0,
    step=0.1,
    format="%0.1f")

age = st.number_input(
    "**Возраст, лет**",
    min_value=0.0,
    step=0.1,
    format="%0.1f")

def calc(weight, height, age):
    if weight == 0 or height == 0 or age == 0:
        st.error("Не заполнены входные данные")
        return

    ubm = 10 * weight + 6.25 * height - 5 * age
    if gender == ":material/man: Мужчина":
        ubm += 5
    else:
        ubm -= 161

    names = ["Уровень Базального Метаболизма (УБМ)"]
    values = [ubm]

    for n, k in [
        ("Сидячая активность", 1.2),
        ("Легкая активность", 1.375),
        ("Средняя активность", 1.55),
		("Полноценная активность", 1.725),
		("Cупер активность", 1.9)
    ]:
        names.append(n)
        values.append(round(ubm * k))

    st.dataframe(
        pd.DataFrame({
            "Показатель": names,
            "Значение, ккал": values,
        }),
        hide_index=True
    )

if st.button("Рассчитать", type="primary"):
    calc(weight, height, age)