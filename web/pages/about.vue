<script setup lang="ts">
import type { IconNode } from 'lucide'
import { Eye, ScanSearch } from 'lucide'
import {
  ABOUT_CRAFTSMANSHIP,
  ABOUT_HERO,
  ABOUT_MISSION,
  ABOUT_STORY,
  ABOUT_SYMBOL,
} from '~/constants/about'

useSeoMeta({
  title: 'About | Oyecouture Bespoke',
  description: ABOUT_HERO.description,
})

const MISSION_ICON: IconNode = ScanSearch
const VISION_ICON: IconNode = Eye
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
      </div>
    </header>

    <!-- 01 / Our Story -->
    <section :id="ABOUT_STORY.id" class="about-section">
      <div class="about-section-body">
        <div class="about-copy">
          <p class="about-section-label">{{ ABOUT_STORY.label }}</p>
          <h2>{{ ABOUT_STORY.heading }}</h2>

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

        <figure class="about-image-panel">
          <img
            src="https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto,w_1000/oyecouture/campaign-4"
            alt="Oyecouture traditional tailoring"
            loading="lazy"
          >
        </figure>
      </div>
    </section>

    <!-- 02 / The Oyecouture Symbol -->
    <section :id="ABOUT_SYMBOL.id" class="about-section is-inverted">
      <div class="about-section-body">
        <div class="about-symbol-mark" aria-hidden="true">
          <span class="about-symbol-ring">
            <img src="/images/brand/oyecouture-icon-gold.png" alt="" loading="lazy" />
          </span>
        </div>

        <div class="about-copy">
          <p class="about-section-label">{{ ABOUT_SYMBOL.label }}</p>
          <h2>{{ ABOUT_SYMBOL.heading }}</h2>

          <div class="about-prose">
            <p v-for="paragraph in ABOUT_SYMBOL.paragraphs" :key="paragraph">{{ paragraph }}</p>
          </div>

          <ul class="about-pillars">
            <li v-for="pillar in ABOUT_SYMBOL.pillars" :key="pillar.title">
              <strong>{{ pillar.value }}</strong>
              <span>{{ pillar.title }}</span>
              <small>{{ pillar.description }}</small>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <!-- 03 / Craftsmanship -->
    <section :id="ABOUT_CRAFTSMANSHIP.id" class="about-section">
      <div class="about-section-body">
        <figure class="about-image-panel is-tailor">
          <img
            src="https://res.cloudinary.com/exgithrw/image/upload/f_auto,q_auto,w_1000/oyecouture/asoke-look-04-3"
            alt="Hand-finished Oyecouture ceremonial attire"
            loading="lazy"
          >
        </figure>

        <div class="about-copy">
          <p class="about-section-label">{{ ABOUT_CRAFTSMANSHIP.label }}</p>
          <h2>{{ ABOUT_CRAFTSMANSHIP.heading }}</h2>

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
      </div>
    </section>

    <!-- 04 / Mission & Vision -->
    <section :id="ABOUT_MISSION.id" class="about-section about-mission-section">
      <p class="about-section-label">{{ ABOUT_MISSION.label }}</p>
      <h2>{{ ABOUT_MISSION.heading }}</h2>

      <div class="about-statements">
        <article
          v-for="statement in ABOUT_MISSION.statements"
          :key="statement.title"
          :class="{ 'is-dark': statement.title === 'Vision' }"
        >
          <div class="about-statement-heading">
            <span class="about-statement-icon" aria-hidden="true">
              <StorefrontIcon
                :icon-node="statement.title === 'Vision' ? VISION_ICON : MISSION_ICON"
                :size="20"
                :stroke-width="1.6"
              />
            </span>
            <h3>{{ statement.title }}</h3>
          </div>
          <p>{{ statement.description }}</p>
        </article>
      </div>
    </section>
  </div>
</template>
