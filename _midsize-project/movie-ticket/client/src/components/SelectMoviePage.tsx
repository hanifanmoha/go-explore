import { useEffect, useState } from 'react'
import useCart from '../hooks/useCart'
import useStep from '../hooks/useStep'
import { MovieCard } from './MovieCard'
import { getApiBaseUrl } from '../env'

export default function SelectMoviePage() {

  const nextStep = useStep((state) => state.nextStep)
  const setMovie = useCart((state) => state.setMovie)

  const [movies, setMovies] = useState([])

  useEffect(() => {
    fetch(`${getApiBaseUrl()}/movies`)
      .then((response) => response.json())
      .then((data) => {
        setMovies(data)
      })
      .catch((error) => {
        console.error('Error fetching movies:', error)
      })

  }, [])

  function handleSelectMovie(movieId: string) {
    const movie = movies.find((movie) => movie.id === movieId)
    if (!movie) return
    setMovie(movie.id, movie.title)
    nextStep()
  }

  return (
    <div className="flex flex-col gap-4">
      {movies.map((movie, index) => (
        <MovieCard key={index} id={movie.id} title={movie.title} description={movie.description} imageUrl={movie.image_url} onSelect={handleSelectMovie} />
      ))}
    </div>
  )
}