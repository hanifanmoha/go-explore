import CheckoutPage from "./components/CheckoutPage"
import Header from "./components/Header"
import SelectMoviePage from "./components/SelectMoviePage"
import SelectSeatPage from "./components/SelectSeatPage"
import SuccessPage from "./components/SuccessPage"
import useStep from "./hooks/useStep"

function App() {

  const step = useStep((state) => state.step)

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
