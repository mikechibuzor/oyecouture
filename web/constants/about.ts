// Copy transcribed verbatim from the approved About design. Values that read as
// uppercase on screen are stored sentence-case and uppercased in CSS, matching
// how the design itself is built.
export const ABOUT_HERO = {
  eyebrow: 'About',
  headingLines: ['Your clothing speaks', 'Long before you do'],
  description:
    "A premium bespoke house for men who understand that presence isn't announced — it's worn. This is confidence, simplicity, and class, made into clothing.",
} as const

export const ABOUT_SECTION_IDS = {
  story: 'our-story',
  symbol: 'symbol',
  craftsmanship: 'craftsmanship',
  mission: 'mission',
} as const

export const ABOUT_TABS = [
  { id: ABOUT_SECTION_IDS.story, label: 'Our Story' },
  { id: ABOUT_SECTION_IDS.symbol, label: 'The Oyecouture Symbol' },
  { id: ABOUT_SECTION_IDS.craftsmanship, label: 'Craftsmanship' },
  { id: ABOUT_SECTION_IDS.mission, label: 'Mission & Vision' },
] as const

export const ABOUT_STORY = {
  id: ABOUT_SECTION_IDS.story,
  label: '01 — Our Story',
  heading: 'A language of confidence, worn into every room',
  paragraphs: [
    'Oyecouture Bespoke was founded on a simple truth: your clothing speaks long before you do. We set out to build a brand for men who understand that presence isn’t announced — it’s worn.',
    'From our earliest pieces, the goal was never just clothing, but a language of confidence, simplicity, and class that any man could carry into a room. Today, Oyecouture stands as a premium entity for two kinds of men — those who already lead with style, and those ready to level up.',
  ],
  stats: [
    { value: 'Two', label: 'Kinds of men' },
    { value: 'One', label: 'Standard — premium' },
  ],
} as const

export const ABOUT_SYMBOL = {
  id: ABOUT_SECTION_IDS.symbol,
  label: '02 — The Oyecouture Symbol',
  heading: 'A philosophy made visible',
  paragraphs: [
    "Our icon is more than a mark. Formed by connecting the letters of our name into one continuous shape, it reads as a pie chart of what we stand for. Look closer and you'll find a window built into the design — our window to exclusivity, a passage into the world of Oyecouture. Simple enough to be replicated anywhere, it remains the most visible signature of who we are.",
  ],
  pillars: [
    { title: 'Confidence', description: 'Presence you wear' },
    { title: 'Simplicity', description: 'Nothing wasted' },
    { title: 'Class', description: 'Quietly assured' },
    { title: 'You', description: 'The window in' },
  ],
} as const

export const ABOUT_CRAFTSMANSHIP = {
  id: ABOUT_SECTION_IDS.craftsmanship,
  label: '03 — Craftsmanship',
  heading: 'The patience of a Ruler, the imagination of a Creator',
  paragraphs: [
    'Every Oyecouture piece is made with the patience of a Ruler and the imagination of a Creator. We are dedicated to timeless quality — pieces, services, and experiences built to outlast trend cycles.',
    'Nothing leaves our hands until it reflects the standard we set for ourselves: premium, deliberate, and unmistakably exclusive.',
  ],
  principles: [
    { number: '01', description: 'Timeless quality, built to outlast trends' },
    { number: '02', description: 'Deliberate, considered from cloth to finish' },
    { number: '03', description: 'Premium and unmistakably exclusive' },
  ],
} as const

export const ABOUT_MISSION = {
  id: ABOUT_SECTION_IDS.mission,
  label: '04 — Mission & Vision',
  heading: 'What we are building',
  statements: [
    {
      title: 'Mission',
      description:
        'To dress the modern man in confidence, simplicity, and class — turning clothing into a language that speaks before a single word is said.',
    },
    {
      title: 'Vision',
      description:
        'To be the definitive symbol of exclusivity in menswear — a brand that connects style, ambition, and identity into one iconic mark, recognized anywhere, worn by those who lead.',
    },
  ],
} as const

// The About page runs its own version of the bespoke band, with a different
// heading and an extra line of copy.
export const ABOUT_BESPOKE_BAND = {
  heading: 'Step into the world of Oyecouture',
  description:
    'A private consultation, hand-selected cloth, and a fit finished over three fittings. Be exclusive.',
} as const
