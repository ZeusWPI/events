function EventsIcon({ className }: { className?: string }) {
  const thickStripes = [
    { x1: 509.21, y1: 123.58, x2: 509.21, y2: 171.94 },
    { x1: 509.93, y1: 849.18, x2: 509.93, y2: 897.55 },
    { x1: 848.4, y1: 516.06, x2: 896.76, y2: 516.06 },
    { x1: 123.1, y1: 516.06, x2: 171.46, y2: 516.06 },
  ];

  const thinStripes = [
    { x1: 178.56, y1: 709.2, x2: 202.95, y2: 695.12 },
    { x1: 820.18, y1: 338.76, x2: 844.58, y2: 324.68 },
    { x1: 317.71, y1: 849.86, x2: 331.79, y2: 825.46 },
    { x1: 695.93, y1: 213.07, x2: 710.01, y2: 188.68 },
    { x1: 710.01, y1: 843.24, x2: 695.93, y2: 818.84 },
    { x1: 298.31, y1: 226.92, x2: 284.22, y2: 202.52 },
    { x1: 839.18, y1: 709.2, x2: 814.78, y2: 695.12 },
    { x1: 183.66, y1: 359.41, x2: 159.27, y2: 345.33 },
  ];

  const sharedProps = {
    fill: "none",
    stroke: "currentColor",
    strokeMiterlimit: 10,
    strokeLinecap: "round" as const,
  };

  return (
    <svg
      viewBox="0 0 1024 1024"
      className={className}
      xmlns="http://www.w3.org/2000/svg"
    >
      <g id="Letter-2">
        <path
          {...sharedProps}
          strokeWidth={66}
          d="M171.46,188.68C256.99,98.54,377.93,42.34,512,42.34s251.12,54.42,336.4,142.03"
        />
        <path
          {...sharedProps}
          strokeWidth={66}
          d="M171.46,837.66c85.53,90.13,206.47,146.34,340.54,146.34s251.12-54.42,336.4-142.03"
        />
        <path
          {...sharedProps}
          strokeLinecap="butt"
          strokeWidth={66}
          d="M188.73,854.96L829.68,166.17,188.73,854.96Z"
        />
        <circle fill="currentColor" cx="509.21" cy="510.56" r="68.12" />
      </g>

      <g id="Clock_Stripes">
        {thickStripes.map((coords, i) => (
          <line key={`thick-${i}`} {...sharedProps} strokeWidth={30} {...coords} />
        ))}
        {thinStripes.map((coords, i) => (
          <line key={`thin-${i}`} {...sharedProps} strokeWidth={20} {...coords} />
        ))}
      </g>
    </svg>
  );
}

export default EventsIcon;
