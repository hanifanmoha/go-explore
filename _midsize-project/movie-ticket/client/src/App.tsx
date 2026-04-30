import { useEffect, useState } from "react"
import CheckoutPage from "./components/CheckoutPage"
import Header from "./components/Header"
import SelectMoviePage from "./components/SelectMoviePage"
import SelectSeatPage from "./components/SelectSeatPage"
import SuccessPage from "./components/SuccessPage"
import useStep from "./hooks/useStep"
import { useLocalStorage } from "./hooks/useLocalStorage"
import { generateRandomID } from "./utils/idGenerator"
import useCart from "./hooks/useCart"
import { getApiBaseUrl } from "./env"

function App() {

  const [userID] = useLocalStorage("userID", generateRandomID(10))
  const step = useStep((state) => state.step)
  const setStep = useStep((state) => state.setStep)
  const setUserID = useCart((state) => state.setUserID)

  useEffect(() => {
    setUserID(userID)
  }, [userID, setUserID])

  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch(`${getApiBaseUrl()}/bookings/${userID}`)
      .then((res) => res.json())
      .then((data) => {
        if (data && Array.isArray(data)) {
          setStep(4)
        } else {
          setStep(1)
        }
        setLoading(false)
      })
      .catch((err) => {
        console.error("Error fetching bookings:", err)
        setLoading(false)
      })
  }, [userID])

  if (loading) {
    return <div className="container mx-auto max-w-xl p-4">
      <Header step={step} />
      <div className="text-center mt-10">
        <p className="text-lg">Loading...</p>
      </div>
    </div>
  }

  return (
    <div className="container mx-auto max-w-xl p-4">
      <Header step={step} />

      {step === 1 && <SelectMoviePage />}
      {step === 2 && <SelectSeatPage />}
      {step === 3 && <CheckoutPage />}
      {step === 4 && <SuccessPage />}
    </div>
  )
}

export default App
