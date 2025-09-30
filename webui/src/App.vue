<script setup lang="ts">
import { RouterLink, RouterView } from 'vue-router'
</script>

<template>
  <header>
    <nav class="navbar navbar-expand-lg bg-body-tertiary">
      <div class="container-fluid">
        <RouterLink to="/" class="navbar-brand">MyHealth</RouterLink>
        <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
          <span class="navbar-toggler-icon"></span>
        </button>
        <div class="collapse navbar-collapse" id="navbarSupportedContent">
          <ul class="navbar-nav me-auto mb-2 mb-lg-0">            
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                Еда
              </a>
              <ul class="dropdown-menu">
                <li><RouterLink to="/food/journal" class="dropdown-item">Журнал приема пищи</RouterLink></li>
                <li><RouterLink to="/food" class="dropdown-item">Управление едой</RouterLink></li>
                <li><RouterLink to="/food/stats" class="dropdown-item">Статистика</RouterLink></li>
              </ul>
            </li>
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                Вес
              </a>
              <ul class="dropdown-menu">
                <li><RouterLink to="/weight" class="dropdown-item">Вес тела</RouterLink></li>
                <li><RouterLink to="/weight/stats" class="dropdown-item">Статистика</RouterLink></li>
              </ul>
            </li>
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                Активность
              </a>
              <ul class="dropdown-menu">
                <li><RouterLink to="/activity/journal" class="dropdown-item">Журнал активности</RouterLink></li>
                <li><RouterLink to="/activity" class="dropdown-item">Управление спортом</RouterLink></li>
                <li><RouterLink to="/activity/stats" class="dropdown-item">Статистика</RouterLink></li>
              </ul>
            </li>
            <li class="nav-item dropdown">
              <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown" aria-expanded="false">
                Настройки
              </a>
              <ul class="dropdown-menu">
                <li><RouterLink to="/settings/calccal" class="dropdown-item">Расчет лимита ккал.</RouterLink></li>
                <li><RouterLink to="/settings/user" class="dropdown-item">Пользователь</RouterLink></li>
              </ul>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </header>

  <div class="container">
    <RouterView />
  </div>
</template>

