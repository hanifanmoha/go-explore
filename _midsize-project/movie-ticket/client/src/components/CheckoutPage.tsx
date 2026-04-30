import { useEffect, useState } from "react"
import useCart from "../hooks/useCart"
import useStep from "../hooks/useStep"
import { MovieCard } from "./MovieCard"
import { generateRandomID } from "../utils/idGenerator"
import { getApiBaseUrl } from "../env"

export default function CheckoutPage() {

  const userID = useCart((state) => state.userID)
  const movieID = useCart((state) => state.movieID)
  const seatID = useCart((state) => state.seatID)
  const seatLabel = useCart((state) => state.seatLabel)
  const prevStep = useStep((state) => state.prevStep)
  const nextStep = useStep((state) => state.nextStep)

  const [movie, setMovie] = useState(null)

  useEffect(() => {
    if (!movieID) return
    fetch(`${getApiBaseUrl()}/movies/${movieID}`)
      .then((response) => response.json())
      .then((data) => {
        setMovie(data)
      })
      .catch((error) => {
        console.error('Error fetching movie:', error)
      })
  }, [movieID])

  function checkoutInfo() {
    return (
      <div>
        <p><strong>Seat:</strong> {`${seatLabel} (ID: ${seatID})`}</p>
      </div>
    )
  }

  async function handleCheckout() {
    try {
      const response = await fetch(`${getApiBaseUrl()}/bookings`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          user_id: userID, // Replace with actual user ID
          movie_id: movieID,
          seat_id: seatID,
        }),
      })

      if (!response.ok) {
        throw new Error("Failed to create booking")
      }

      nextStep()
    } catch (error) {
      alert("Error creating booking: " + error.message)
    }
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <MovieCard
        id={movie?.id || ""}
        title={movie?.title || ""}
        description={checkoutInfo()}
        imageUrl={movie?.image_url || ""}
      />
      <div className="flex flex-col gap-4 w-full">
        <button className="btn btn-primary w-full" onClick={handleCheckout}>Checkout</button>
        <button className="btn w-full" onClick={prevStep}>Back</button>
      </div>
    </div>
  )
}