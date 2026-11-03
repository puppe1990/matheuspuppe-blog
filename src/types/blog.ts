export interface Author {
  id: string;
  name: string;
  role: string;
  bio: string;
  avatar: string;
  twitter?: string;
  github?: string;
  articleCount?: number;
}

export type ContentBlockType = 
  | 'paragraph'
  | 'heading2'
  | 'heading3'
  | 'quote'
  | 'code'
  | 'callout'
  | 'takeaways'
  | 'interactive-benchmark';

export interface ContentBlock {
  type: ContentBlockType;
  text?: string;
  code?: string;
  language?: string;
  caption?: string;
  authorQuote?: string;
  items?: string[];
  benchmarkData?: {
    title: string;
    description: string;
    metrics: { label: string; valueA: number; valueB: number; unit: string }[];
  };
}

export interface Article {
  id: string;
  slug: string;
  title: string;
  subtitle: string;
  leadParagraph: string;
  category: 'IA & Modelos' | 'Engenharia de Sistemas' | 'Hardware & Chips' | 'Design & UX' | 'Segurança & Criptografia' | 'Opinião & Futuro';
  readTimeMinutes: number;
  publishedAt: string; // ISO date
  coverImage: string;
  author: Author;
  tags: string[];
  blocks: ContentBlock[];
  claps: number;
  views: number;
  featured?: boolean;
  tier?: 1 | 2 | 3; // Tier 1 = lead, Tier 2 = secondary featured, Tier 3 = standard
}

export interface Comment {
  id: string;
  articleId: string;
  authorName: string;
  authorRole: string;
  content: string;
  createdAt: string;
  claps: number;
}

export interface NewsletterSubscription {
  email: string;
  subscribedAt: string;
}
