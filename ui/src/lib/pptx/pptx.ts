import { convertEventToModel, Event } from "../types/event";

import { format } from "date-fns";
import { nlBE } from "date-fns/locale";
import pptxgen from "pptxgenjs";
import QrCodeWithLogo from "qrcode-with-logos";
import { apiGet } from "../api/query";
import { capitalize } from "../utils/utils";

const masterName = "ZEUS_WPI_TEMPLATE"

const colorZeus = "#FF7F00"
const colorWhite = "#FFFFFF"
const colorBlack = "#000000"

const urlWebsite = "https://zeus.gent"
const urlInstagram = "https://www.instagram.com/zeuswpi/"

export async function generatePptx(events: Event[]) {
  const pptx = new pptxgen();

  generateMaster(pptx)
  generateIntro(pptx)

  for (const event of events) {
    await generateEventSlide(pptx, event)
  }

  await generateOutro(pptx)

  return pptx.writeFile({ fileName: "zeus_events.pptx" })
}

async function generateMaster(pptx: pptxgen) {
  pptx.defineSlideMaster({
    title: masterName,
    background: { color: colorWhite },
    objects: [
      // Bottom rectangle
      { rect: { x: 0, y: "90%", w: "100%", h: "10%", fill: { color: colorZeus } } },
      // Website icon + text
      {
        image: {
          path: "/images/website.png",
          x: "4%",
          y: "92%",
          w: 0.35,
          h: 0.35,
        }
      },
      {
        text: {
          text: urlWebsite,
          options: { x: "7.5%", y: "95%", w: "86%", align: "left", color: colorWhite }
        }
      },
      // Instagram text + icon
      {
        text: {
          text: "@zeuswpi",
          options: { x: "4%", y: "95%", w: "92%", align: "right", color: colorWhite }
        }
      },
      {
        image: {
          path: "/images/instagram.png",
          x: "78%",
          y: "92%",
          w: 0.35,
          h: 0.35,
          hyperlink: {
            tooltip: "Zeus instagram",
            url: urlInstagram,
          },
        }
      },
      // Top line
      {
        line: { x: 0, y: "15%", w: "75%", h: 0, line: { color: colorZeus, width: 2 } },
      },
      // Bottom line
      {
        line: { x: "63%", y: "85%", w: "38%", h: 0, line: { color: colorZeus, width: 2 } },
      },
    ],
  })
}

function generateIntro(pptx: pptxgen) {
  const slide = pptx.addSlide()

  slide.background = { "color": colorZeus }

  // Main logo
  slide.addImage({
    path: "images/zeus_white.svg",
    x: "20%",
    y: "20%",
    w: 6,
    h: 3.733,
    hyperlink: {
      tooltip: "Zeus website",
      url: "https://zeus.gent",
    },
    shadow: {
      type: "outer",
      blur: 5,
      offset: 5,
      opacity: 0.5,
    }
  })

  // Mascot
  slide.addImage({
    path: "images/mascot.png",
    x: "67.7%",
    y: "45%",
  })
}

async function generateOutro(pptx: pptxgen) {
  const slide = pptx.addSlide()

  slide.background = { "color": colorZeus }

  // Website
  slide.addImage({
    path: "/images/website.png",
    x: "21.5%",
    y: "7%",
    w: 0.75,
    h: 0.75,
  })

  slide.addText(
    urlWebsite,
    { x: "8.5%", y: "25%", w: "50%", fontSize: 27, color: colorWhite }
  )

  const qrWebsite = await getQrCode(urlWebsite, colorBlack, colorZeus)
  slide.addImage({
    path: qrWebsite,
    x: "10.5%",
    y: "33%",
    h: 3,
    w: 3,
  })

  // Instagram
  slide.addImage({
    path: "/images/instagram.png",
    x: "68%",
    y: "7%",
    w: 0.75,
    h: 0.75,
  })

  slide.addText(
    "@zeuswpi",
    { x: "61.5%", y: "25%", w: "50%", fontSize: 27, color: colorWhite }
  )

  const qrInstagram = await getQrCode(urlInstagram, colorBlack, colorZeus)
  slide.addImage({
    path: qrInstagram,
    x: "57%",
    y: "33%",
    h: 3,
    w: 3,
  })
}

async function generateEventSlide(pptx: pptxgen, event: Event) {
  const slide = pptx.addSlide({ masterName })

  // Title
  slide.addText(
    event.name,
    { x: 0.0, y: "8%", w: "100%", align: "center", fontSize: 27, color: colorBlack },
  )

  // General information
  // Date
  slide.addImage({
    path: "images/calendar.png",
    x: "7.5%",
    y: "22%",
    w: 0.35,
    h: 0.35,
  })
  slide.addText(
    capitalize(formatDate(event.startTime)),
    { x: "10%", y: "25%", w: "40%", fontSize: 18, color: colorBlack },
  )

  // Time
  slide.addImage({
    path: "images/clock.png",
    x: "7.2%",
    y: "27.4%",
    w: 0.42,
    h: 0.42,
  })
  slide.addText(
    formatTime(event.startTime),
    { x: "10%", y: "31%", w: "40%", fontSize: 18, color: colorBlack },
  )

  // Location
  slide.addImage({
    path: "images/location.png",
    x: "7.6%",
    y: "38.5%",
    w: 0.35,
    h: 0.35,
  })
  slide.addText(
    event.location,
    { x: "10%", y: "38%", h: "40%", w: "40%", valign: "top", fontSize: 18, color: colorBlack },
  )

  const populated = await getEvent(event)

  if (populated.posters.some(p => !p.scc)) {
    // Event has a big poster
    await generateEventSlidePoster(slide, populated)
  } else {
    await generateEventSlideNoPoster(slide, populated)
  }
}

async function generateEventSlidePoster(slide: pptxgen.Slide, event: Event) {
  const poster = `/api/poster/${event.posters.find(p => !p.scc)?.id ?? 0}?original=true`

  if (poster) {
    slide.addImage({
      path: poster,
      x: "63%",
      y: "20%",
      h: 3.5,
      w: 2.4748,
    })

  }

  // Qr code to event
  const qr = await getQrCode(event.url)

  slide.addImage({
    path: qr,
    x: "40%",
    y: "53.5%",
    h: 2,
    w: 2,
  })
}

async function generateEventSlideNoPoster(slide: pptxgen.Slide, event: Event) {
  const qr = await getQrCode(event.url)

  slide.addImage({
    path: qr,
    x: "62.2%",
    y: "25%",
    h: 3.3,
    w: 3.3,
  })
}

function formatDate(date: Date) {
  return format(date, "eeee dd MMMM", { locale: nlBE })
}

function formatTime(date: Date) {
  return format(date, "HH'h'mm")
}

async function getEvent({ id }: Pick<Event, "id">) {
  return (await apiGet(`event/${id}`, convertEventToModel)).data
}

async function getQrCode(url: string, cornerColor: string = colorZeus, background: string = colorWhite): Promise<string> {
  const qr = new QrCodeWithLogo({
    content: url,
    width: 640,
    nodeQrCodeOptions: {
      color: {
        light: background,
      },
    },
    logo: {
      src: "/images/mascot.png",
    },
    dotsOptions: {
      type: "fluid",
    },
    cornersOptions: {
      type: "rounded",
      color: cornerColor,
    },
  })

  const img = await qr.getImage()
  return img.src
}

