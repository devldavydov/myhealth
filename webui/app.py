import streamlit as st

# Index
index_page = st.Page(
    "index.py",
    title="Главная",
    default=True,
    icon=":material/home:"
)

# Food
food_page = st.Page(
    "food/food.py",
    url_path="food",
    title="Управление едой",
    icon=":material/settings:"
)
food_journal_page = st.Page(
    "food/journal.py",
    url_path="food-journal",
    title="Журнал приёма пищи",
    icon=":material/fork_spoon:"
)
food_stats_page = st.Page(
    "food/stats.py",
    url_path="food-stats",
    title="Статистика",
    icon=":material/bar_chart:"
)

# Weight
weight_page = st.Page(
    "weight/weight.py",
    url_path="weight",
    title="Вес тела",
    icon=":material/scale:"
)
weight_stats_page = st.Page(
    "weight/stats.py",
    url_path="weight-stats",
    title="Статистика",
    icon=":material/bar_chart:"
)

# Activity
activity_page = st.Page(
    "activity/sport.py",
    url_path="activity",
    title="Управление спортом",
    icon=":material/settings:"
)
activity_journal_page = st.Page(
    "activity/journal.py",
    url_path="activity-journal",
    title="Журнал активности",
    icon=":material/exercise:"
)
activity_stats_page = st.Page(
    "activity/stats.py",
    url_path="activity-stats",
    title="Статистика",
    icon=":material/bar_chart:"
)

# Settings
calc_cal_page = st.Page(
    "settings/calc_cal.py",
    url_path="settings-calc-cal",
    title="Расчет лимита ккал",
    icon=":material/function:"
)
user_settings_page = st.Page(
    "settings/user.py",
    url_path="settings-user",
    title="Пользователь",
    icon=":material/account_circle:"
)

pg = st.navigation(
    {
        "": [index_page],
        "Еда": [food_journal_page, food_page, food_stats_page],
        "Активность": [activity_journal_page, activity_page, activity_stats_page],
        "Вес": [weight_page, weight_stats_page],
        "Настройки": [calc_cal_page, user_settings_page]
    },
    expanded=True
)

pg.run()