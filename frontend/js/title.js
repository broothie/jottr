
export default function setSubtitle(subtitle) {
  if (subtitle) {
    document.title = `jottr - ${subtitle}`
  } else {
    document.title = 'jottr'
  }
}
