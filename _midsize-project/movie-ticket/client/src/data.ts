import type { Movie, SeatData } from "./model"

export const movies: Movie[] = [
  {
    id: "the-bourne-identity",
    title: 'The Bourne Identity',
    description: 'A man is picked up by a fishing boat, bullet-riddled and suffering from amnesia, before racing to elude assassins and attempting to regain his memory.',
    imageUrl: 'https://m.media-amazon.com/images/M/MV5BYTk1ZTcyMWMtMWUxYS00MmEzLTlmODYtOTk1MGRjOTg1ZjlmXkEyXkFqcGc@._V1_FMjpg_UX1000_.jpg',
  },
  {
    id: "the-bourne-supremacy",
    title: 'The Bourne Supremacy',
    description: 'When Jason Bourne is framed for a CIA operation gone awry, he is forced to resume his former life as a trained assassin to survive.',
    imageUrl: 'https://m.media-amazon.com/images/M/MV5BZTU4ZDgyYjgtODA0Mi00MmE3LTgzYWQtZjc1YTFiMTczZTQ3XkEyXkFqcGc@._V1_FMjpg_UX1014_.jpg',
  },
  {
    id: "the-bourne-ultimatum",
    title: 'The Bourne Ultimatum',
    description: 'Jason Bourne dodges a ruthless C.I.A. official and his Agents from a new assassination program while searching for the origins of his life as a trained killer.',
    imageUrl: 'https://m.media-amazon.com/images/M/MV5BYzE3ZGM4MzctZjU5MC00NWE2LTg5ZjYtMDFiM2ZlMWQ1MjkwXkEyXkFqcGc@._V1_FMjpg_UY3000_.jpg',
  },
]

export const seats: SeatData[][] = Array(4).fill(Array(8).fill(false)).map((row, rowIndex) => {
  return row.map((_, colIndex) => {
    return {
      id: `${rowIndex + 1}_${colIndex + 1}`,
      isAvailable: (rowIndex * 8 + colIndex) % 5 !== 0
    }
  })
})