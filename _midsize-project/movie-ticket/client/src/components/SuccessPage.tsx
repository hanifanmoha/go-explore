import { useEffect, useState } from "react"
import { useLocalStorage } from "../hooks/useLocalStorage"
import { generateRandomID } from "../utils/idGenerator"

export default function SuccessPage() {

  const [bookings, setBookings] = useState([])
  const [userID, setUserID] = useLocalStorage("userID", "")

  useEffect(() => {
    fetch(`http://localhost:6001/bookings/${userID}`)
      .then((res) => res.json())
      .then((data) => {
        if (data && Array.isArray(data)) {
          setBookings(data)
        }
      })
      .catch((err) => {
        console.error("Error fetching bookings:", err)
      })
  }, [userID])

  function bookOtherTicket() {
    setUserID(generateRandomID(10))
    window.location.reload()
  }

  let booking = null
  if (bookings.length > 0) {
    booking = bookings[0]
  }

  return (
    <div className="flex flex-col items-center gap-4 h-100 justify-center">
      <h1 className="text-2xl font-bold">Success!</h1>
      <p>Your movie ticket has been booked successfully.</p>
      <button className="btn btn-primary" onClick={() => bookOtherTicket()}>Book Another Ticket</button>
    </div>
  )
}