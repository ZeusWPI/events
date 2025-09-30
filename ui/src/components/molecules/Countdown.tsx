import { useEffect, useState } from "react";
import NumberFlow, { NumberFlowGroup } from "@number-flow/react";

interface Props {
  goalDate: Date;
}

export function Countdown({ goalDate }: Props) {
  const [timeRemaining, setTimeRemaining] = useState(goalDate.getTime() - new Date().getTime())

  useEffect(() => {
    if (timeRemaining <= 0) return

    const interval = setInterval(() => {
      setTimeRemaining(prev => Math.max(prev - 1000, 0))
    }, 1000)

    return () => clearInterval(interval)
  }, [timeRemaining])

  if (timeRemaining <= 0) {
    return <p className="font-bold text-red-500">!</p>
  }

  const seconds = Math.floor((timeRemaining / 1000) % 60);
  const minutes = Math.floor((timeRemaining / (1000 * 60)) % 60);
  const hours = Math.floor((timeRemaining / (1000 * 60 * 60)) % 24);
  const days = Math.floor(timeRemaining / (1000 * 60 * 60 * 24));

  return (
    <NumberFlowGroup>
      <div className="flex">
        <NumberFlow
          trend={-1}
          value={days}
          format={{ minimumIntegerDigits: 2 }}
        />
        <p>:</p>
        <NumberFlow
          trend={-1}
          value={hours}
          digits={{ 1: { max: 5 } }}
          format={{ minimumIntegerDigits: 2 }}
        />
        <p>:</p>
        <NumberFlow
          trend={-1}
          value={minutes}
          digits={{ 1: { max: 5 } }}
          format={{ minimumIntegerDigits: 2 }}
        />
        <p>:</p>
        <NumberFlow
          trend={-1}
          value={seconds}
          digits={{ 1: { max: 5 } }}
          format={{ minimumIntegerDigits: 2 }}
        />
      </div>
    </NumberFlowGroup>
  )
}
