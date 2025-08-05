interface Props {
  title?: string;
  description?: string;
}

export const NoItems = ({ title, description }: Props) => {
  return (
    <div className="flex flex-col items-center space-y-4 pt-48">
      <h3 className="text-lg font-semibold">{title}</h3>
      <h5 className="text-md text-muted-foreground">{description}</h5>
    </div>
  )
}
