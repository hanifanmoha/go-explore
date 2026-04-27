import { movies } from '../data'
import useCart from '../hooks/useCart'
import useStep from '../hooks/useStep'
import { MovieCard } from './MovieCard'

export default function SelectMoviePage() {

  const nextStep = useStep((state) => state.nextStep)
  const setMovie = useCart((state) => state.setMovie)

  function handleSelectMovie(movieId: string) {
    const movie = movies.find((movie) => movie.id === movieId)
    if (!movie) return
    setMovie(movie.id, movie.title)
    nextStep()
  }

  return (
    <div className="flex flex-col gap-4">
      {movies.map((movie, index) => (
        <MovieCard key={index} id={movie.id} title={movie.title} description={movie.description} imageUrl={movie.imageUrl} onSelect={handleSelectMovie} />
      ))}
    </div>
  )
}