interface MovieCardProps {
  id: string
  title: string
  description: string | React.ReactNode
  imageUrl: string
  onSelect?: (movieId: string) => void
}

export function MovieCard({ id, title, description, imageUrl, onSelect }: MovieCardProps) {
  return (
    <div className="card bg-base-100 w-full shadow-sm">
      <figure>
        <img
          className='max-w-1/2'
          src={imageUrl}
          alt={title} />
      </figure>
      <div className="card-body">
        <h2 className="card-title">{title}</h2>
        <div>{description}</div>
        <div className="card-actions justify-end">
          {onSelect && <button className="btn btn-primary" onClick={() => onSelect(id)}>Select Seat</button>}
        </div>
      </div>
    </div>
  )
}