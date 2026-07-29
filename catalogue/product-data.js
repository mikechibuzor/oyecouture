// Shared product catalog for Oyecouture. Used by the product detail page.
// Slot naming convention (matches all listing pages): 'oye-p-<id>-<1|2|3>'
export const products = [
  // ---- Suits ----
  { id: 'monarch', name: 'The Monarch Three-Piece', cat: 'Three-Piece', section: 'suits', price: 850000, tag: 'Bespoke', hasSizes: false,
    description: 'A three-piece cut for boardrooms and altars alike — structured shoulder, tapered trouser and a waistcoat that holds its line all day. Made to a single measure in your choice of cloth.' },
  { id: 'sovereign', name: 'The Sovereign Suit', cat: 'Two-Piece', section: 'suits', price: 620000, tag: '', hasSizes: false,
    description: 'Our most requested suit: a clean two-button silhouette in a mid-weight wool blend, built for the man who wears one suit exceptionally well.' },
  { id: 'regent', name: 'The Regent Two-Piece', cat: 'Two-Piece', section: 'suits', price: 540000, tag: '', hasSizes: false,
    description: "A softer shoulder and lighter cloth make the Regent the suit you reach for when the room calls for less armor, more ease." },
  { id: 'heir', name: 'The Heir Blazer', cat: 'Blazers', section: 'suits', price: 380000, tag: 'New', hasSizes: false,
    description: 'An unstructured blazer that moves like a shirt and holds like a jacket — worn open over knitwear, or buttoned for meetings that matter.' },
  { id: 'ambassador', name: 'The Ambassador Three-Piece', cat: 'Three-Piece', section: 'suits', price: 720000, tag: 'New', hasSizes: false,
    description: "Peak lapels and a fuller drape give the Ambassador its presence — reserved for days you're representing more than yourself." },
  { id: 'chairman', name: 'The Chairman Three-Piece', cat: 'Three-Piece', section: 'suits', price: 680000, tag: '', hasSizes: false,
    description: 'Double-breasted and unapologetically formal, cut for the head of the table.' },
  { id: 'diplomat', name: 'The Diplomat Two-Piece', cat: 'Two-Piece', section: 'suits', price: 490000, tag: '', hasSizes: false,
    description: 'A travel-friendly two-piece in a fine, breathable wool — the creases fall out by the time you land.' },
  { id: 'weekend', name: 'The Weekend Linen Blazer', cat: 'Blazers', section: 'suits', price: 275000, tag: '', hasSizes: false,
    description: 'Washed linen, patch pockets, no lining — the blazer built to survive a Lagos afternoon.' },

  // ---- Trad Attires ----
  { id: 'agbada', name: 'The Royal Agbada', cat: 'Agbada', section: 'trad', price: 950000, tag: 'Bespoke', hasSizes: false,
    description: 'Our most elaborate piece: hand-embroidered agbada in richly woven aso-oke, built over three fittings for the day everyone remembers.' },
  { id: 'heritage', name: 'Heritage Aso-Oke Set', cat: 'Aso-Oke', section: 'trad', price: 720000, tag: '', hasSizes: false,
    description: 'A complete aso-oke set — buba, sokoto and fila — woven in a heritage pattern reworked for a modern fit.' },
  { id: 'fila', name: 'Aso-Oke Fila & Wrapper', cat: 'Aso-Oke', section: 'trad', price: 410000, tag: 'New', hasSizes: false,
    description: 'Fila and wrapper in matching aso-oke cloth — the finishing pieces for any traditional outfit.' },
  { id: 'emir', name: 'The Emir Kaftan', cat: 'Kaftans', section: 'trad', price: 290000, tag: '', hasSizes: false,
    description: 'Hand-embroidered collar and cuffs on a fluid, breathable kaftan — Friday prayers to family gatherings.' },
  { id: 'senator', name: 'Embroidered Senator Kaftan', cat: 'Kaftans', section: 'trad', price: 240000, tag: '', hasSizes: false,
    description: 'The senator cut, sharpened: slim through the body, finished with subtle thread embroidery at the neckline.' },
  { id: 'classic', name: 'Classic White Kaftan', cat: 'Kaftans', section: 'trad', price: 165000, tag: '', hasSizes: false,
    description: 'Pure white cotton, no embellishment — the kaftan that goes with everything you already own.' },
  { id: 'ceremonial', name: 'The Ceremonial Agbada', cat: 'Agbada', section: 'trad', price: 780000, tag: '', hasSizes: false,
    description: 'A lighter agbada for daytime ceremonies, in breathable cotton brocade with understated embroidery.' },
  { id: 'executive', name: 'The Executive Agbada', cat: 'Agbada', section: 'trad', price: 690000, tag: 'New', hasSizes: false,
    description: 'Streamlined agbada for the man who wants the silhouette without the volume — tailored, not draped.' },
  { id: 'wedding', name: 'Aso-Oke Wedding Ensemble', cat: 'Aso-Oke', section: 'trad', price: 860000, tag: 'Bespoke', hasSizes: false,
    description: "A complete aso-oke ensemble for the groom's party, woven to match across every member of your train." },

  // ---- Casuals · Men ----
  { id: 'linen', name: 'The Weekend Linen Shirt', cat: 'Shirts', section: 'casuals-men', price: 78000, tag: '', hasSizes: true,
    description: 'Pure linen, cut generously through the body — the shirt that wants a slow Saturday.' },
  { id: 'oxford', name: 'The Oxford Overshirt', cat: 'Shirts', section: 'casuals-men', price: 86000, tag: 'New', hasSizes: true,
    description: 'A heavier cotton overshirt that layers over a tee or under a blazer without complaint.' },
  { id: 'polo', name: 'The Resort Polo', cat: 'Shirts', section: 'casuals-men', price: 54000, tag: '', hasSizes: true,
    description: 'Piqué cotton polo with a longer body and clean placket — resort dressing without trying too hard.' },
  { id: 'chino', name: 'Tailored Chino Trouser', cat: 'Trousers', section: 'casuals-men', price: 92000, tag: 'New', hasSizes: true,
    description: 'A tailored taper through the leg, finished with a soft, worn-in hand — dresses up or down.' },
  { id: 'jogger', name: 'The Utility Jogger', cat: 'Trousers', section: 'casuals-men', price: 74000, tag: '', hasSizes: true,
    description: 'Tapered jogger in a structured cotton twill, built to read tailored even at rest.' },
  { id: 'short', name: 'Tailored Cotton Short', cat: 'Trousers', section: 'casuals-men', price: 58000, tag: '', hasSizes: true,
    description: 'A tailored cotton short cut like a trouser — the piece that makes shorts look intentional.' },
  { id: 'crew', name: 'The Crewneck Sweater', cat: 'Knitwear', section: 'casuals-men', price: 88000, tag: '', hasSizes: true,
    description: 'Midweight cotton crewneck, fully fashioned for a clean line under a blazer or on its own.' },
  { id: 'cardigan', name: 'The Merino Cardigan', cat: 'Knitwear', section: 'casuals-men', price: 96000, tag: 'New', hasSizes: true,
    description: 'Merino wool cardigan with horn buttons — the layer that finishes a look without adding bulk.' },

  // ---- Casuals · Women ----
  { id: 'wrap', name: 'The Silk Wrap Blouse', cat: 'Tops', section: 'casuals-women', price: 82000, tag: 'New', hasSizes: true,
    description: 'A silk wrap blouse that ties at the waist — one piece, several necklines, depending on the day.' },
  { id: 'linenshirt', name: 'The Linen Shirt Dress', cat: 'Tops', section: 'casuals-women', price: 96000, tag: '', hasSizes: true,
    description: 'A linen shirt dress cut long enough to wear alone, or open over trousers.' },
];

export const sections = {
  'suits': { label: 'Suits', href: 'Oyecouture Suits.dc.html' },
  'trad': { label: 'Trad Attires', href: 'Oyecouture Trad Attires.dc.html' },
  'casuals-men': { label: 'Casuals · Men', href: 'Oyecouture Casuals Men.dc.html' },
  'casuals-women': { label: 'Casuals · Women', href: 'Oyecouture Casuals Women.dc.html' },
};

export const sizeOptions = ['S', 'M', 'L', 'XL'];

export function findProduct(id) {
  return products.find((p) => p.id === id) || null;
}

export function fmt(n) {
  return '₦' + n.toLocaleString('en-US');
}
