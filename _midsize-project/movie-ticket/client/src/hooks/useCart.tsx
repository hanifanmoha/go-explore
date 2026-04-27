import { create } from 'zustand'

interface CartState {
  seatID: string,
  movieID: string,
  movieTitle: string,
  setSeat: (seatID: string) => void,
  setMovie: (movieID: string, movieTitle: string) => void,
}

const useCart = create<CartState>((set) => ({
  seatID: "",
  movieID: "",
  movieTitle: "",
  setSeat: (seatID: string) => set({ seatID }),
  setMovie: (movieID: string, movieTitle: string) => set({ movieID, movieTitle }),
}))

export default useCart;