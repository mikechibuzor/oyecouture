<script setup lang="ts">
import type { IconNode } from 'lucide'
import {
  ArrowRight,
  CalendarDays,
  Check,
  Mail,
  MessageCircle,
  Star,
} from 'lucide'
import { buildWhatsAppUrl } from '~/utils/whatsapp'

interface ContactInfoItem {
  readonly label: string
  readonly value: string
  readonly href?: string
}

interface ContactSocialLink {
  readonly name: string
  readonly handle: string
  readonly href: string
  readonly icon: IconNode
}

const { data: siteSettings } = await useSiteSettings()
const selectedTopic = ref('General')
const selectedRating = ref(0)
const senderName = ref('')
const senderEmail = ref('')
const feedbackMessage = ref('')
const feedbackSent = ref(false)
const feedbackError = ref('')
const isFeedbackSubmitting = ref(false)
const { submitContactFeedback } = useContactFeedback()
const { appointmentWhatsAppUrl, opensWhatsApp } = await useAppointmentWhatsApp()

const configuredWhatsAppNumber = computed(() =>
  String(siteSettings.value?.whatsappNumber ?? '').replace(/\D/g, ''),
)
const displayedPhoneNumber = computed(() => `+${configuredWhatsAppNumber.value}`)
const whatsappHref = computed(() =>
  buildWhatsAppUrl(configuredWhatsAppNumber.value) || undefined,
)

const contactInfo = computed<readonly ContactInfoItem[]>(() => [
  { label: 'Email', value: 'hello@oyecouturebespoke.com', href: 'mailto:hello@oyecouturebespoke.com' },
  ...(whatsappHref.value
    ? [{ label: 'Phone & WhatsApp', value: displayedPhoneNumber.value, href: whatsappHref.value }]
    : []),
])

const CONTACT_TOPICS = ['Fitting', 'Order', 'Feedback', 'Press'] as const
const contactSocials = computed<readonly ContactSocialLink[]>(() => [
  ...(whatsappHref.value
    ? [{
        name: 'WhatsApp',
        handle: displayedPhoneNumber.value,
        href: whatsappHref.value,
        icon: MessageCircle,
      }]
    : []),
  { name: 'Email', handle: 'hello@oyecouturebespoke.com', href: 'mailto:hello@oyecouturebespoke.com', icon: Mail },
  {
    name: 'Oyecouture',
    handle: 'Book a private fitting',
    href: appointmentWhatsAppUrl.value,
    icon: CalendarDays,
  },
])

const stars = computed(() =>
  Array.from({ length: 5 }, (_, starIndex) => ({
    value: starIndex + 1,
    label: `${starIndex + 1} star${starIndex === 0 ? '' : 's'}`,
    active: selectedRating.value >= starIndex + 1,
  })),
)

const submitFeedback = async (): Promise<void> => {
  if (!senderName.value.trim() || !senderEmail.value.trim() || !feedbackMessage.value.trim()) {
    feedbackError.value = 'Please add your name, email and message.'
    return
  }

  feedbackError.value = ''
  isFeedbackSubmitting.value = true

  try {
    await submitContactFeedback({
      name: senderName.value,
      email: senderEmail.value,
      topic: selectedTopic.value,
      rating: selectedRating.value,
      message: feedbackMessage.value,
    })
    feedbackSent.value = true
  }
  catch {
    feedbackError.value = 'We could not send this yet. Please try again.'
  }
  finally {
    isFeedbackSubmitting.value = false
  }
}

const resetFeedback = (): void => {
  senderName.value = ''
  senderEmail.value = ''
  feedbackMessage.value = ''
  selectedTopic.value = 'General'
  selectedRating.value = 0
  feedbackSent.value = false
  feedbackError.value = ''
}

useStorefrontSeo({
  title: 'Contact',
  description:
    "Questions, feedback, a fitting to book, or simply a word about how we're doing. Reach Oyecouture directly.",
})
</script>

<template>
  <div class="contact-page">
    <header class="contact-hero">
      <div class="contact-hero-content">
        <p class="eyebrow">Contact</p>
        <h1>
          We are listening,<br>
          always
        </h1>
        <p>
          Questions, feedback, a fitting to book, or simply a word about how we're doing.
          reach Oyecouture directly, or find us in every space we keep.
        </p>
      </div>
    </header>

    <section id="feedback" class="contact-section">
      <div class="contact-split">
        <div class="contact-reach">
          <p class="contact-section-label">01 / Reach Us</p>
          <h2>Speak to the house directly</h2>
          <p>
            Every message is read by someone who cares about the cut. Expect a
            considered reply within one business day.
          </p>

          <dl class="contact-info-grid">
            <div v-for="infoItem in contactInfo" :key="infoItem.label">
              <dt>{{ infoItem.label }}</dt>
              <dd>
                <a v-if="infoItem.href" :href="infoItem.href" target="_blank" rel="noopener">
                  {{ infoItem.value }}
                </a>
                <span v-else>{{ infoItem.value }}</span>
              </dd>
            </div>
          </dl>

          <a
            class="contact-primary-link"
            :href="appointmentWhatsAppUrl"
            :target="opensWhatsApp ? '_blank' : undefined"
            :rel="opensWhatsApp ? 'noopener' : undefined"
          >
            <span>Book a fitting</span>
            <StorefrontIcon :icon-node="ArrowRight" :size="14" :stroke-width="2" aria-hidden="true" />
          </a>
        </div>

        <form class="contact-feedback-form" @submit.prevent="submitFeedback">
          <div v-if="feedbackSent" class="contact-feedback-success">
            <span aria-hidden="true">
              <StorefrontIcon :icon-node="Check" :size="34" :stroke-width="1.6" />
            </span>
            <h3>Thank you, {{ senderName }}</h3>
            <p>
              Your message has reached Oyecouture. We read every note personally
              and will reply within one business day.
            </p>
            <button type="button" @click="resetFeedback">Send another</button>
          </div>

          <div v-else>
            <p class="contact-section-label">02 / Share Feedback</p>
            <h3>Tell us how we're doing</h3>

            <div class="contact-form-row">
              <label>
                <span>Name</span>
                <input v-model="senderName" type="text" placeholder="Your name">
              </label>
              <label>
                <span>Email</span>
                <input v-model="senderEmail" type="email" placeholder="you@email.com">
              </label>
            </div>

            <fieldset class="contact-topic-group">
              <legend>Reason for reaching out</legend>
              <button
                v-for="topic in CONTACT_TOPICS"
                :key="topic"
                type="button"
                :class="{ 'is-active': selectedTopic === topic }"
                @click="selectedTopic = topic"
              >
                {{ topic }}
              </button>
            </fieldset>

            <fieldset class="contact-rating-group">
              <legend>How was your experience?</legend>
              <button
                v-for="star in stars"
                :key="star.value"
                type="button"
                :aria-label="star.label"
                :class="{ 'is-active': star.active }"
                @click="selectedRating = star.value"
              >
                <StorefrontIcon :icon-node="Star" :size="30" :stroke-width="1.3" aria-hidden="true" />
              </button>
            </fieldset>

            <label class="contact-message-field">
              <span>Your message</span>
              <textarea
                v-model="feedbackMessage"
                placeholder="Share your thoughts, a compliment, or how we can do better..."
              />
            </label>

            <p v-if="feedbackError" class="contact-form-error">{{ feedbackError }}</p>

            <button class="contact-submit-button" type="submit" :disabled="isFeedbackSubmitting">
              {{ isFeedbackSubmitting ? 'Sending...' : 'Send feedback' }}
            </button>
            <p class="contact-form-note">We'll never share your details. Read by a human, always.</p>
          </div>
        </form>
      </div>
    </section>

    <section id="social" class="contact-social-section">
      <div class="contact-social-content">
        <p class="contact-section-label">03 / Contact Oyecouture</p>
        <h2>Choose how to reach us</h2>
        <p>
          Start a conversation, ask about an order or arrange a private fitting.
        </p>

        <div class="contact-social-grid">
          <a
            v-for="socialLink in contactSocials"
            :key="socialLink.name"
            :href="socialLink.href"
            :target="socialLink.href.startsWith('http') ? '_blank' : undefined"
            :rel="socialLink.href.startsWith('http') ? 'noopener' : undefined"
          >
            <span class="contact-social-icon" aria-hidden="true">
              <StorefrontIcon :icon-node="socialLink.icon" :size="22" :stroke-width="1.5" />
            </span>
            <span>
              <strong>{{ socialLink.name }}</strong>
              <small>{{ socialLink.handle }}</small>
            </span>
            <StorefrontIcon
              class="contact-social-arrow"
              :icon-node="ArrowRight"
              :size="18"
              :stroke-width="1.6"
              aria-hidden="true"
            />
          </a>
        </div>
      </div>
    </section>

    <section class="contact-atelier-section">
      <div>
        <p class="contact-section-label">Prefer to speak in person</p>
        <h2>Visit Oyecouture</h2>
        <p>
          Book a private consultation, feel the cloth, and let us take your measure.
          Be exclusive.
        </p>
        <a
          :href="appointmentWhatsAppUrl"
          :target="opensWhatsApp ? '_blank' : undefined"
          :rel="opensWhatsApp ? 'noopener' : undefined"
        >Book an appointment</a>
      </div>
    </section>
  </div>
</template>
