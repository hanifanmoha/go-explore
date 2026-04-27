import { movies } from "../data"
import useCart from "../hooks/useCart"
import useStep from "../hooks/useStep"
import { MovieCard } from "./MovieCard"

export default function CheckoutPage() {

  const movieID = useCart((state) => state.movieID)
  const seatID = useCart((state) => state.seatID)
  const prevStep = useStep((state) => state.prevStep)
  const nextStep = useStep((state) => state.nextStep)

  const movie = movies.find((movie) => movie.id === movieID)

  function checkoutInfo() {
    return (
      <div>
        <p><strong>Seat:</strong> {seatID}</p>
      </div>
    )
  }

  function handleCheckout() {
    nextStep()
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <MovieCard
        id={movie?.id || ""}
        title={movie?.title || ""}
        description={checkoutInfo()}
        imageUrl={movie?.imageUrl || ""}
      />
      <div className="flex flex-col gap-4 w-full">
        <button className="btn btn-primary w-full" onClick={handleCheckout}>Checkout</button>
        <button className="btn w-full" onClick={prevStep}>Back</button>
      </div>
    </div>
  )
}