import test from 'node:test';
import assert from 'node:assert/strict';
import {
  filterArticles,
  calculateReadingTime,
  estimateWordCount,
  toggleBookmark,
  calculateClap,
  formatPublicationDate,
  calculateScrollProgress,
  calculateRemainingReadingTime,
  getCategoryCounts,
  buildArticleShareData,
} from './blogLogic';
import { Article } from '../types/blog';

const mockAuthor = {
  id: 'author-1',
  name: 'Helena Vance',
  role: 'Pesquisadora Chefe de Sistemas',
  bio: 'Especialista em compiladores e sistemas distribuídos.',
  avatar: '/avatar1.jpg',
};

const mockArticles: Article[] = [
  {
    id: 'art-1',
    slug: 'artigo-ia',
    title: 'A Nova Fronteira dos Agentes Autônomos',
    subtitle: 'Por que o raciocínio superou a escala pura',
    leadParagraph: 'Modelos de linguagem estão migrando de geração estocástica para planejamento multi-passos.',
    category: 'IA & Modelos',
    readTimeMinutes: 7,
    publishedAt: '2026-03-15T10:00:00Z',
    coverImage: '/img1.jpg',
    author: mockAuthor,
    tags: ['IA', 'Agentes', 'LLM'],
    blocks: [
      { type: 'paragraph', text: 'Primeiro parágrafo com dez palavras para teste de contagem.' },
      { type: 'takeaways', items: ['Ponto um', 'Ponto dois'] },
    ],
    claps: 120,
    views: 1400,
  },
  {
    id: 'art-2',
    slug: 'artigo-rust',
    title: 'Migrando de Go para Rust em 2026',
    subtitle: 'Zero garbage collection em serviços de baixa latência',
    leadParagraph: 'Como eliminamos spikes de p99 em infraestrutura de pagamentos.',
    category: 'Engenharia de Sistemas',
    readTimeMinutes: 4,
    publishedAt: '2026-02-10T14:30:00Z',
    coverImage: '/img2.jpg',
    author: mockAuthor,
    tags: ['Rust', 'Go', 'Backend'],
    blocks: [{ type: 'paragraph', text: 'Texto do segundo artigo sobre Rust.' }],
    claps: 85,
    views: 950,
  },
];

test('filterArticles filters by search query matching title', () => {
  const result = filterArticles(mockArticles, 'Agentes');
  assert.equal(result.length, 1);
  assert.equal(result[0].id, 'art-1');
});

test('filterArticles filters by category', () => {
  const result = filterArticles(mockArticles, '', 'Engenharia de Sistemas');
  assert.equal(result.length, 1);
  assert.equal(result[0].id, 'art-2');
});

test('filterArticles filters by maxReadTime', () => {
  const result = filterArticles(mockArticles, '', 'Todos', 5);
  assert.equal(result.length, 1);
  assert.equal(result[0].id, 'art-2');
});

test('filterArticles filters by bookmark IDs', () => {
  const result = filterArticles(mockArticles, '', 'Todos', null, true, ['art-1']);
  assert.equal(result.length, 1);
  assert.equal(result[0].id, 'art-1');
});

test('calculateReadingTime calculates minutes correctly', () => {
  assert.equal(calculateReadingTime(200, 200), 1);
  assert.equal(calculateReadingTime(201, 200), 2);
  assert.equal(calculateReadingTime(0, 200), 1);
});

test('estimateWordCount calculates words in blocks', () => {
  const words = estimateWordCount(mockArticles[0].blocks);
  // 'Primeiro parágrafo com dez palavras para teste de contagem.' -> 9 words
  // 'Ponto um' -> 2 words, 'Ponto dois' -> 2 words = 13 words
  assert.equal(words, 13);
});

test('toggleBookmark adds and removes id', () => {
  const initial: string[] = ['art-1'];
  const added = toggleBookmark(initial, 'art-2');
  assert.deepEqual(added, ['art-1', 'art-2']);
  const removed = toggleBookmark(added, 'art-1');
  assert.deepEqual(removed, ['art-2']);
});

test('calculateClap increments properly', () => {
  assert.equal(calculateClap(10, 1), 11);
  assert.equal(calculateClap(10, 5), 15);
});

test('calculateScrollProgress calculates scroll percentage with bounds 0 to 100', () => {
  assert.equal(calculateScrollProgress(0, 1000), 0);
  assert.equal(calculateScrollProgress(500, 1000), 50);
  assert.equal(calculateScrollProgress(1000, 1000), 100);
  assert.equal(calculateScrollProgress(1200, 1000), 100);
  assert.equal(calculateScrollProgress(-50, 1000), 0);
  assert.equal(calculateScrollProgress(100, 0), 0);
});

test('calculateRemainingReadingTime computes remaining reading minutes accurately', () => {
  assert.equal(calculateRemainingReadingTime(10, 0), 10);
  assert.equal(calculateRemainingReadingTime(10, 50), 5);
  assert.equal(calculateRemainingReadingTime(10, 90), 1);
  assert.equal(calculateRemainingReadingTime(10, 100), 0);
});

test('getCategoryCounts aggregates counts including Todos', () => {
  const counts = getCategoryCounts(mockArticles);
  assert.equal(counts['Todos'], 2);
  assert.equal(counts['IA & Modelos'], 1);
  assert.equal(counts['Engenharia de Sistemas'], 1);
});

test('buildArticleShareData generates compliant Web Share API payload', () => {
  const payload = buildArticleShareData(
    { title: 'Artigo Tech', subtitle: 'Dissecação de código' },
    'https://sintaxe.tech/artigo-tech'
  );
  assert.equal(payload.title, 'Artigo Tech');
  assert.equal(payload.text, 'Artigo Tech — Dissecação de código');
  assert.equal(payload.url, 'https://sintaxe.tech/artigo-tech');
});
