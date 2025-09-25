import streamlit as st

def create_food_pages():
    return [
        st.Page(
            "food/food.py",
            url_path="food",
            title="Управление едой",
            icon=":material/settings:"
        ),
        st.Page(
            "food/journal.py",
            url_path="food-journal",
            title="Журнал приёма пищи",
            icon=":material/fork_spoon:"
        ),
        st.Page(
            "food/stats.py",
            url_path="food-stats",
            title="Статистика",
            icon=":material/bar_chart:"
        ),
        st.Page(
            "food/edit.py",
            url_path="food-edit",
            title="Редактирование еды",
            icon=":material/edit:"
        ),
        st.Page(
            "food/create.py",
            url_path="food-create",
            title="Создание еды",
            icon=":material/add:"
        )
    ]

def create_weight_pages():
    return [
        st.Page(
            "weight/weight.py",
            url_path="weight",
            title="Вес тела",
            icon=":material/scale:"
        ),
        st.Page(
            "weight/stats.py",
            url_path="weight-stats",
            title="Статистика",
            icon=":material/bar_chart:"
        )
    ]

def create_activity_pages():
    return [
        st.Page(
            "activity/sport.py",
            url_path="activity",
            title="Управление спортом",
            icon=":material/settings:"
        ),
        st.Page(
            "activity/journal.py",
            url_path="activity-journal",
            title="Журнал активности",
            icon=":material/exercise:"
        ),
        st.Page(
            "activity/stats.py",
            url_path="activity-stats",
            title="Статистика",
            icon=":material/bar_chart:"
        )
    ]

def create_settings_pages():
    return [
        st.Page(
            "settings/calc_cal.py",
            url_path="settings-calc-cal",
            title="Расчет лимита ккал",
            icon=":material/function:"
        ),
        st.Page(
            "settings/user.py",
            url_path="settings-user",
            title="Пользователь",
            icon=":material/account_circle:"
        )
    ]

def create_navigation():
    # Create page router
    pages = [
        st.Page(
            "index.py",
            title="Главная",
            default=True,
            icon=":material/home:"
        )
    ]

    pages.extend(create_food_pages())
    pages.extend(create_weight_pages())
    pages.extend(create_activity_pages())
    pages.extend(create_settings_pages())

    pg = st.navigation(pages, expanded=True, position="hidden")

    # Add sidebar links
    st.sidebar.markdown("##### Еда")
    st.sidebar.page_link(
        "food/journal.py",
        label="Журнал приёма пищи",
        icon=":material/fork_spoon:"
    )
    st.sidebar.page_link(
        "food/food.py",
        label="Управление едой",
        icon=":material/settings:"
    )
    st.sidebar.page_link(
        "food/stats.py",
        label="Статистика",
        icon=":material/bar_chart:"
    )

    st.sidebar.markdown("##### Вес")
    st.sidebar.page_link(
        "weight/weight.py",
        label="Вес тела",
        icon=":material/scale:"
    )
    st.sidebar.page_link(
        "weight/stats.py",
        label="Статистика",
        icon=":material/bar_chart:"
    )

    st.sidebar.markdown("##### Активность")
    st.sidebar.page_link(
        "activity/journal.py",
        label="Журнал активности",
        icon=":material/exercise:"
    )
    st.sidebar.page_link(
        "activity/sport.py",
        label="Управление спортом",
        icon=":material/settings:"
    )    
    st.sidebar.page_link(
        "activity/stats.py",
        label="Статистика",
        icon=":material/bar_chart:"
    )

    st.sidebar.markdown("##### Настройки")
    st.sidebar.page_link(
        "settings/calc_cal.py",
        label="Расчет лимита ккал",
        icon=":material/function:"
    )
    st.sidebar.page_link(
        "settings/user.py",
        label="Пользователь",
        icon=":material/account_circle:"
    ) 

    return pg

def main():
    # Global settings
    st.set_page_config(layout="wide")
    st.session_state["api_url"] = "http://127.0.0.1:8080/api"

    # Run
    pg = create_navigation()
    pg.run()

if __name__ == "__main__":
    main()