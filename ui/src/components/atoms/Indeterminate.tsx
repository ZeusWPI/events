export function Indeterminate() {
  return (
    <div className="flex justify-center items-center w-full h-[80%]">
      <div className="h-1.5 w-[90%] bg-secondary-foreground rounded-full overflow-hidden">
        <div className="progress h-full bg-primary rounded-full left-right"></div>
      </div>
    </div>
  );
}
