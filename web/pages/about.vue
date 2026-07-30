<script setup lang="ts">
import {
  ABOUT_CRAFTSMANSHIP,
  ABOUT_HERO,
  ABOUT_MISSION,
  ABOUT_STORY,
  ABOUT_SYMBOL,
  ABOUT_TABS,
} from '~/constants/about'

useSeoMeta({
  title: 'About — Oyecouture Bespoke',
  description: ABOUT_HERO.description,
})
</script>

<template>
  <div class="about-page">
    <header class="about-hero">
      <div class="about-hero-content">
        <p class="eyebrow">{{ ABOUT_HERO.eyebrow }}</p>
        <h1>
          <span v-for="line in ABOUT_HERO.headingLines" :key="line">{{ line }}</span>
        </h1>
        <p class="about-hero-description">{{ ABOUT_HERO.description }}</p>

        <nav class="about-tabs" aria-label="About sections">
          <a v-for="tab in ABOUT_TABS" :key="tab.id" :href="`#${tab.id}`">
            {{ tab.label }}
          </a>
        </nav>
      </div>
    </header>

    <!-- 01 — Our Story -->
    <section :id="ABOUT_STORY.id" class="about-section">
      <p class="about-section-label">{{ ABOUT_STORY.label }}</p>
      <h2>{{ ABOUT_STORY.heading }}</h2>

      <div class="about-section-body">
        <div class="about-prose">
          <p v-for="paragraph in ABOUT_STORY.paragraphs" :key="paragraph">{{ paragraph }}</p>
        </div>

        <dl class="about-stats">
          <div v-for="stat in ABOUT_STORY.stats" :key="stat.label">
            <dt>{{ stat.value }}</dt>
            <dd>{{ stat.label }}</dd>
          </div>
        </dl>
      </div>
    </section>

    <!-- 02 — The Oyecouture Symbol -->
    <section :id="ABOUT_SYMBOL.id" class="about-section is-inverted">
      <p class="about-section-label">{{ ABOUT_SYMBOL.label }}</p>
      <h2>{{ ABOUT_SYMBOL.heading }}</h2>

      <div class="about-section-body">
        <div class="about-prose">
          <p v-for="paragraph in ABOUT_SYMBOL.paragraphs" :key="paragraph">{{ paragraph }}</p>
        </div>

        <ul class="about-pillars">
          <li v-for="pillar in ABOUT_SYMBOL.pillars" :key="pillar.title">
            <strong>{{ pillar.title }}</strong>
            <span>{{ pillar.description }}</span>
          </li>
        </ul>
      </div>
    </section>

    <!-- 03 — Craftsmanship -->
    <section :id="ABOUT_CRAFTSMANSHIP.id" class="about-section">
      <p class="about-section-label">{{ ABOUT_CRAFTSMANSHIP.label }}</p>
      <h2>{{ ABOUT_CRAFTSMANSHIP.heading }}</h2>

      <div class="about-section-body">
        <div class="about-prose">
          <p v-for="paragraph in ABOUT_CRAFTSMANSHIP.paragraphs" :key="paragraph">
            {{ paragraph }}
          </p>
        </div>

        <ol class="about-principles">
          <li v-for="principle in ABOUT_CRAFTSMANSHIP.principles" :key="principle.number">
            <span class="about-principle-number">{{ principle.number }}</span>
            <span>{{ principle.description }}</span>
          </li>
        </ol>
      </div>
    </section>

    <!-- 04 — Mission & Vision -->
    <section :id="ABOUT_MISSION.id" class="about-section">
      <p class="about-section-label">{{ ABOUT_MISSION.label }}</p>
      <h2>{{ ABOUT_MISSION.heading }}</h2>

      <div class="about-statements">
        <article v-for="statement in ABOUT_MISSION.statements" :key="statement.title">
          <h3>{{ statement.title }}</h3>
          <p>{{ statement.description }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
