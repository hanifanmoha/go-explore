import { seats } from "../data"
import useCart from "../hooks/useCart"
import useStep from "../hooks/useStep"

interface SeatProps {
  isSelected: boolean
  isAvailable: boolean
  onSelect: () => void
}

function Seat({ isSelected, isAvailable, onSelect }: SeatProps) {

  let bgColor = "bg-gray-300"
  if (isSelected) bgColor = "bg-primary"
  if (!isAvailable) bgColor = "bg-neutral"

  function handleSelect() {
    if (!isAvailable) return
    onSelect()
  }

  return (
    <div className={`w-12 h-12 ${bgColor} rounded-sm cursor-pointer hover:opacity-75`} onClick={handleSelect} />
  )
}

interface SeatData {
  id: string
  isAvailable: boolean
}

interface SeatsProps {
  seats: SeatData[][],
  seatID: string | null,
  setSeat: (seatID: string) => void,
}

function Seats({ seats, seatID, setSeat }: SeatsProps) {

  return (
    <div className="flex flex-col gap-2 justify-center items-center">
      {seats.map((row, rowIndex) => (
        <div key={rowIndex} className="flex gap-2">
          {row.map((seat, seatIndex) => (
            <Seat
              key={seatIndex}
              isSelected={seat.id === seatID}
              isAvailable={seat.isAvailable}
              onSelect={() => setSeat(seat.id)}
            />
          ))}
        </div>
      ))}
    </div>
  )
}

function Legend() {
  return <div className="flex gap-4 justify-center items-center">
    <div className="flex gap-2 items-center">
      <div className="w-4 h-4 bg-gray-300 rounded-sm" />
      <span>Available</span>
    </div>
    <div className="flex gap-2 items-center">
      <div className="w-4 h-4 bg-primary rounded-sm" />
      <span>Selected</span>
    </div>
    <div className="flex gap-2 items-center">
      <div className="w-4 h-4 bg-neutral rounded-sm" />
      <span>Unavailable</span>
    </div>
  </div>
}

export default function SelectSeatPage() {

  const movieTitle = useCart((state) => state.movieTitle)
  const seatID = useCart((state) => state.seatID)
  const setSeat = useCart((state) => state.setSeat)

  const nextStep = useStep((state) => state.nextStep)
  const prevStep = useStep((state) => state.prevStep)

  return (
    <div className="flex flex-col items-center gap-4">
      <h1 className="text-2xl font-bold">{movieTitle}</h1>
      <Seats seats={seats} seatID={seatID} setSeat={setSeat} />
      <Legend />
      <div className="flex flex-col gap-4 mt-12 w-full">
        <button className="btn btn-primary w-full" onClick={nextStep}>Continue</button>
        <button className="btn w-full" onClick={prevStep}>Back</button>
      </div>
    </div>
  )
}