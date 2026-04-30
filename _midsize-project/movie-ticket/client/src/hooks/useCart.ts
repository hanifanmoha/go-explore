import { create } from 'zustand'
import { generateRandomID } from '../utils/idGenerator';

interface CartState {
  userID: string,
  seatID: string,
  seatLabel: string,
  movieID: string,
  movieTitle: string,
  setUserID: (userID: string) => void,
  setSeat: (seatID: string, seatLabel: string) => void,
  setMovie: (movieID: string, movieTitle: string) => void,
}

const useCart = create<CartState>((set) => ({
  userID: "",
  seatID: "",
  seatLabel: "",
  movieID: "",
  movieTitle: "",
  setUserID: (userID: string) => set({ userID }),
  setSeat: (seatID: string, seatLabel: string) => set({ seatID, seatLabel }),
  setMovie: (movieID: string, movieTitle: string) => set({ movieID, movieTitle }),
}))

export default useCart;