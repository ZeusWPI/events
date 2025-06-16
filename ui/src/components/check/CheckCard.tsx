import { Check, CheckSource } from "@/lib/types/check"
import { Check as CheckIcon, ChevronDownIcon, ChevronUpIcon, ClipboardCheckIcon, ClipboardXIcon, XIcon } from "lucide-react"
import { useState } from "react"
import { Button } from "../ui/button"
import { Tooltip, TooltipContent, TooltipTrigger } from "../ui/tooltip"
import { useCheckToggle } from "@/lib/api/check"
import { toast } from "sonner"

type Props = {
  check: Check
}

export function CheckCard({ check }: Props) {
  const [showError, setShowError] = useState(false)
  const [toggleStatus, setToggleStatus] = useState(false)

  console.log(check)

  const checkToggle = useCheckToggle()

  const toggleShowError = () => setShowError(prev => !prev)
  const toggleDone = () => {
    setToggleStatus(true)

    console.log(check)

    checkToggle.mutate(check, {
      onSuccess: () => toast.success("Success"),
      onError: error => toast.error("Failed", { description: error.message }),
      onSettled: () => setToggleStatus(false)
    })
  }

  return (
    <div className="grid grid-cols-10">
      <div>
        {check.done ? (
          <CheckIcon className="text-green-500" />
        ) : (
          <XIcon className="text-red-500" />
        )}
      </div>
      <div className="col-span-8">
        <span>{check.description}</span>
      </div>
      <div>
        {check.source === CheckSource.Automatic ? (
          check.error && (
            <Button onClick={toggleShowError} size="icon" variant="outline">
              <Tooltip>
                <TooltipTrigger asChild>
                  {showError ? <ChevronUpIcon /> : <ChevronDownIcon />}
                </TooltipTrigger>
                <TooltipContent>
                  <span>{showError ? "Hide error" : "Show error"}</span>
                </TooltipContent>
              </Tooltip>
            </Button>
          )
        ) : (
          <div className="flex justify-end">
            <Button onClick={toggleDone} size="icon" variant="outline" disabled={toggleStatus}>
              <Tooltip>
                <TooltipTrigger asChild>
                  {check.done ? <ClipboardXIcon /> : <ClipboardCheckIcon />}
                </TooltipTrigger>
                <TooltipContent>
                  <span>{check.done ? "Mark as undone" : "Mark as done"}</span>
                </TooltipContent>
              </Tooltip>
            </Button>
          </div>
        )}
      </div>
      <div>
        {showError && (
          <span className="text-red-500">{check.error}</span>
        )}
      </div>
    </div>
  )
}
