import { Nav } from './components/Nav'
import { Hero } from './components/Hero'
import { SocialProof } from './components/SocialProof'
import { Features } from './components/Features'
import { Metrics } from './components/Metrics'
import { FinalCta } from './components/FinalCta'
import { Footer } from './components/Footer'

export function LandingPage() {
  return (
    <div className="flex min-h-svh flex-col bg-white font-dm-sans">
      <Nav />
      <Hero />
      <SocialProof />
      <Features />
      <Metrics />
      <FinalCta />
      <Footer />
    </div>
  )
}
