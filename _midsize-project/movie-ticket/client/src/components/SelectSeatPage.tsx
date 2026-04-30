import { useEffect, useState } from "react"
import useCart from "../hooks/useCart"
import useStep from "../hooks/useStep"

interface SeatProps {
  isSelected: boolean
  isAvailable: boolean
  onSelect: () => void
  label: string
}

function Seat({ isSelected, isAvailable, label, onSelect }: SeatProps) {

  let bgColor = "bg-gray-300"
  if (isSelected) bgColor = "bg-primary"
  if (!isAvailable) bgColor = "bg-neutral"

  function handleSelect() {
    if (!isAvailable) return
    onSelect()
  }

  return (
    <div className={`w-12 h-12 ${bgColor} rounded-t-lg cursor-pointer hover:opacity-75 flex justify-center items-center`} onClick={handleSelect} >
      {(!isAvailable || isSelected) && <span className="text-xs text-white font-bold">{label}</span>}
      {(isAvailable && !isSelected) && <span className="text-xs text-black font-bold">{label}</span>}
    </div>
  )
}

interface SeatData {
  id: string
  seat_number: string
  status: string
  isAvailable: boolean
}

interface SeatsProps {
  seats: SeatData[][],
  seatID: string | null,
  setSeat: (seatID: string, seatLabel: string) => void,
}

function Seats({ seats, seatID, setSeat }: SeatsProps) {

  return (
    <div className="flex flex-col gap-2 justify-center items-center">
      {seats.map((row, rowIndex) => (
        <div key={rowIndex} className="flex gap-2">
          {row.map((seat, seatIndex) => (
            <Seat
              key={seatIndex}
              label={seat.seat_number}
              isSelected={seat.id === seatID}
              isAvailable={seat.isAvailable}
              onSelect={() => setSeat(seat.id, seat.seat_number)}
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
  const movieID = useCart((state) => state.movieID)
  const seatID = useCart((state) => state.seatID)
  const setSeat = useCart((state) => state.setSeat)

  const nextStep = useStep((state) => state.nextStep)
  const prevStep = useStep((state) => state.prevStep)

  const [seats, setSeats] = useState([])

  useEffect(() => {
    // GET http://localhost:6001/movies/{movieID}/seats
    if (!movieID) return

    fetch(`http://localhost:6001/movies/${movieID}/seats`)
      .then((response) => response.json())
      .then((data) => {
        setSeats(data)
      })
      .catch((error) => {
        console.error('Error fetching seats:', error)
      })
  }, [movieID])

  function get2DSeats() {
    const seatMap: { [key: string]: SeatData[] } = {}
    for (let seat of seats) {
      const row = seat.seat_number[0] // A
      if (!seatMap[row]) seatMap[row] = []
      seatMap[row].push({
        id: seat.id,
        seat_number: seat.seat_number,
        status: seat.status,
        isAvailable: seat.status === 'available'
      })
    }
    const sortedRows = Object.keys(seatMap).sort()
    return sortedRows.map((row) => {
      const sortedSeats = seatMap[row].sort((a, b) => {
        const numA = parseInt(a.seat_number.slice(1))
        const numB = parseInt(b.seat_number.slice(1))
        return numA - numB
      })
      return sortedSeats
    })
  }

  return (
    <div className="flex flex-col items-center gap-4">
      <h1 className="text-2xl font-bold">{movieTitle}</h1>
      <Seats seats={get2DSeats()} seatID={seatID} setSeat={setSeat} />
      <Legend />
      <div className="flex flex-col gap-4 mt-12 w-full">
        <button className="btn btn-primary w-full" onClick={nextStep}>Continue</button>
        <button className="btn w-full" onClick={prevStep}>Back</button>
      </div>
    </div>
  )
}