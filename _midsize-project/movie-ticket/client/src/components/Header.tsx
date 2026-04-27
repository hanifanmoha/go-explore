import { useEffect, useState } from "react";

interface HeaderProps {
  step: number
}

export default function Header({ step }: HeaderProps) {

  const [timeLeft, setTimeLeft] = useState(5 * 60);

  useEffect(() => {
    // Stop at 0
    if (timeLeft <= 0) return;

    // Set interval to update every 1 second
    const timer = setInterval(() => {
      setTimeLeft((prev) => prev - 1);
    }, 1000);

    // Cleanup interval on component unmount
    return () => clearInterval(timer);
  }, [timeLeft]);


  function getStepClass(stepNumber: number) {
    if (stepNumber <= step) {
      return "step-primary"
    }
    return ""
  }

  function getStepIcon(stepNumber: number) {
    if (stepNumber < step) {
      return "✓"
    }
    return ""
  }

  return (
    <div className="mx-auto flex flex-col items-center mb-12 gap-8">
      <ul className="steps gap-4">
        <li data-content={getStepIcon(1)} className={`step ${getStepClass(1)}`}>Select Movie</li>
        <li data-content={getStepIcon(2)} className={`step ${getStepClass(2)}`}>Select Seat</li>
        <li data-content={getStepIcon(3)} className={`step ${getStepClass(3)}`}>Checkout</li>
      </ul>
      {step < 4 && <div className="grid grid-flow-col gap-5 text-center auto-cols-max">
        <div className="flex flex-col p-2 bg-neutral rounded-box text-neutral-content">
          <span className="font-mono text-5xl">
            <span>{Math.floor(timeLeft / 60).toString().padStart(2, "0")}</span>
          </span>
          min
        </div>
        <div className="flex flex-col p-2 bg-neutral rounded-box text-neutral-content">
          <span className="font-mono text-5xl">
            <span>{(timeLeft % 60).toString().padStart(2, "0")}</span>
          </span>
          sec
        </div>
      </div>}
    </div>
  )
}