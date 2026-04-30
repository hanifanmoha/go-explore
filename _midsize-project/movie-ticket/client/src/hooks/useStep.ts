import { create } from 'zustand'

interface StepState {
  step: number
  setStep: (step: number) => void
  nextStep: () => void
  prevStep: () => void
}

const useStep = create<StepState>((set) => ({
  step: 0,
  setStep: (step: number) => set({ step }),
  nextStep: () => set((state) => ({ step: state.step + 1 })),
  prevStep: () => set((state) => ({ step: state.step - 1 })),
}))

export default useStep;