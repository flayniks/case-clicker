import json

# ---------------------------------------------------------------- name tables
BOT_NAMES = [
 # era 1 - scripts and macros
 "Macro Script","Click Bot","Auto Mouse","Rubber Band","Turbo Mouse","Script Kiddie",
 "Keyboard Macro","Mouse Jiggler","Hotkey Daemon","Loop Script",
 # era 2 - accounts and farms
 "Trade Bot","Alt Account","Farm Account","Idle Miner","Bot Lobby","Smurf Army",
 "Key Farmer","Drop Farmer","Case Runner","Inventory Mule",
 # era 3 - networks
 "Botnet","Proxy Swarm","VPN Cluster","Packet Miner","Zombie Fleet","Relay Node",
 "Shadow Net","Grid Worker","Mesh Farm","Dark Pool",
 # era 4 - hardware
 "Server Rack","GPU Rig","Mining Farm","Data Centre","Cooling Array","ASIC Cluster",
 "Blade Server","Fibre Trunk","Rack Row","Cold Storage Vault",
 # era 5 - artificial intelligence
 "Neural Net","AI Trader","Machine Broker","Deep Model","Transformer Farm","Inference Cluster",
 "Sentient Script","Agent Swarm","Model Foundry","Oracle Engine",
 # era 6 - quantum
 "Quantum Clicker","Qubit Array","Entangled Pair","Superposition Farm","Decoherence Engine",
 "Quantum Annealer","Photonic Core","Spin Lattice","Topological Weave","Coherence Reactor",
 # era 7 - fundamental forces
 "Time Dilator","Chrono Loop","Causality Engine","Entropy Reverser","Vacuum Tap",
 "Zero Point Array","Higgs Collector","Gluon Forge","Fusion Spire","Antimatter Well",
 # era 8 - stellar engineering
 "Dyson Swarm","Star Lifter","Neutron Press","Pulsar Array","Black Hole Tap",
 "Accretion Mill","Quasar Engine","Magnetar Coil","Singularity Core","Event Horizon Loom",
 # era 9 - galactic
 "Galactic Exchange","Nebula Refinery","Cluster Broker","Void Harvester","Dark Matter Net",
 "Cosmic Web Node","Filament Tap","Supercluster Mint","Great Attractor","Hubble Engine",
 # era 10 - omniversal
 "Multiverse Market","Timeline Broker","Reality Forge","Axiom Engine","Infinity Loom",
 "Omega Processor","Concept Miner","The Architect","Prime Mover","The Unboxer",
]
BOT_ICONS = ["📜","🤖","🖱️","🪢","⚡","👾","⌨️","🎛️","🔑","🔁",
 "💱","👥","🌾","⛏️","🎮","🕶️","🗝️","📦","🏃","🐴",
 "🕸️","🌐","🔒","📡","🧟","📶","🌑","⚙️","🕳️","🎱",
 "🖥️","🎨","🏭","🏢","❄️","🔲","🗄️","🔌","📚","🧊",
 "🧠","📈","🤝","🔬","🔮","☁️","👁️","🐝","🏗️","🔱",
 "⚛️","🎲","🔗","🌫️","💫","🧮","💡","🧲","🪢","☢️",
 "⏳","🔄","➿","🌡️","🕯️","0️⃣","🧿","🎇","🔥","💥",
 "🛰️","⭐","🌟","📻","🕳️","🌀","💠","🧭","⚫","🎭",
 "🌌","☄️","🔭","🌊","🌑","🕷️","🧵","🏛️","🧲","📐",
 "♾️","⏱️","🛠️","📜","🧬","Ω","💭","🏗️","👑","🎁"]

UPG_NAMES = [
 "Mouse Grease","Mechanical Switches","Gaming Chair","Wrist Brace","Optical Sensor",
 "Braided Cable","1000Hz Polling","Lighter Springs","Gaming Mousepad","Custom Grip Tape",
 "Drop Hack","Macro Firmware","Overclocked Mouse","Titanium Switches","Hall Effect Sensor",
 "Zero Debounce","8K Wireless","Featherweight Shell","Magnetic Glide","Carbon Skates",
 "Global Elite Aim","Muscle Memory","Reflex Training","Pro Wrist","Tournament Form",
 "Frame Perfect","Tick Rate Mastery","Input Buffer","Predictive Aim","Flow State",
 "Neural Link","Cortex Implant","Synaptic Boost","Myelin Overdrive","Motor Cortex Tune",
 "Dopamine Loop","Thalamic Bypass","Reflex Arc Rewire","Brainstem Overclock","Nerve Lattice",
 "Bionic Hand","Servo Tendons","Graphene Bones","Piezo Fingertips","Hydraulic Knuckles",
 "Actuator Array","Exo Glove","Titanium Phalanx","Nanofibre Sinew","Kinetic Amplifier",
 "Time Compression","Bullet Time","Chrono Trigger","Tachyon Input","Retrocausal Click",
 "Temporal Buffer","Clock Skew","Second Splitter","Instant Register","Frozen Moment",
 "Quantum Tunnel","Superposed Input","Probability Bias","Wavefunction Collapse","Entangled Cursor",
 "Zeno Effect","Uncertainty Exploit","Phase Shift","Tunnelling Diode","Observer Effect",
 "Gravity Assist","Inertial Dampener","Mass Driver","Kinetic Singularity","Warp Tap",
 "Fold Space","Graviton Press","Tidal Force","Escape Velocity","Lagrange Hold",
 "Stellar Impact","Supernova Tap","Nova Cascade","Solar Flare","Coronal Surge",
 "Photon Hammer","Plasma Lance","Fusion Strike","Radiance Burst","Helios Protocol",
 "Reality Click","Axiom Override","Concept Tap","Narrative Edit","Root Access",
 "Admin Privilege","Source Patch","Divine Mandate","First Cause","The Final Click",
]
UPG_ICONS = ["🖱️","⌨️","🪑","🩹","🔦","🧵","📊","🌀","🟦","🎗️",
 "📉","💾","🔥","🔩","🧲","🚫","📡","🪶","🧊","⬛",
 "🎖️","🧠","🏃","💪","🏆","🎞️","📶","📥","🎯","🌊",
 "🧠","🔌","⚡","🧬","🗺️","💊","🛤️","🔁","🎚️","🕸️",
 "🦾","🔧","💎","👆","🛠️","🎛️","🧤","🦴","🧶","📣",
 "⏳","🔫","⏰","💫","↩️","📦","🕰️","✂️","⚡","🧊",
 "🌀","🎲","🎰","📉","🔗","🧿","❓","🌈","🔌","👁️",
 "🪐","🛡️","🚀","⚫","🌌","📐","🏋️","🌊","🛸","⚖️",
 "☄️","💥","✨","🔆","🌞","🔨","🗡️","💢","🌟","🛕",
 "🌍","📜","💭","✍️","🔑","👑","🧑‍💻","⛩️","🥚","🏁"]

CASE_NAMES = [
 "Chroma","Prisma","Spectrum","Gamma","Falchion","Shadow","Revolver","Huntsman","Breakout","Vanguard",
 "Phoenix","Winter Offensive","Esports","Bravo","Hydra","Glove","Clutch","Danger Zone","Horizon","Prisma II",
 "Fracture","Snakebite","Dreams","Recoil","Revolution","Kilowatt","Gallery","Fever","Anubis","Riptide",
 "Titanium","Obsidian","Cobalt","Tungsten","Iridium","Osmium","Palladium","Rhodium","Platinum","Meteoric Iron",
 "Emerald","Sapphire","Ruby","Amethyst","Onyx","Opal","Topaz","Garnet","Diamond","Star Sapphire",
 "Quantum","Photon","Graviton","Neutrino","Tachyon","Boson","Lepton","Quark","Gluon","Higgs",
 "Nebula","Pulsar","Quasar","Magnetar","Supernova","Neutron Star","White Dwarf","Red Giant","Protostar","Hypernova",
 "Singularity","Event Horizon","Accretion","Ergosphere","Hawking","Kerr","Schwarzschild","Wormhole","White Hole","Naked Singularity",
 "Galactic","Andromeda","Magellanic","Triangulum","Sombrero","Whirlpool","Pinwheel","Cartwheel","Great Attractor","Laniakea",
 "Multiverse","Timeline","Paradox","Axiom","Infinity","Omega","Genesis","Architect","Prime","Black Market",
]
CASE_ART = ["📦","🎁","🗃️","🧰","💼","🗳️","📮","🧳","🎒","🛍️"]

ERAS = ["Starting Out","Going Wide","Networked","Hardware","Artificial","Quantum",
        "Fundamental","Stellar","Galactic","Omniversal"]
CASE_ERAS = ["Standard","Operation","Modern","Metals","Gemstone","Particle",
             "Stellar","Gravitational","Galactic","Beyond"]

def j(x): return json.dumps(x, ensure_ascii=False)

# ---------------------------------------------------------------- generation
out = []

# ---- BOTS: cost x9.0 per tier, output x9.55 so later bots are better value
out.append("const BOTS=[")
for i,n in enumerate(BOT_NAMES):
    cost = 15 * (7.8 ** i)
    cps  = 0.15 * (8.3 ** i)
    out.append("{id:%s,n:%s,ico:%s,era:%d,base:%s,cps:%s}," % (
        j("b%d"%(i+1)), j(n), j(BOT_ICONS[i]), i//10, repr(cost), repr(cps)))
out.append("];")

# ---- CLICK UPGRADES: flat adds, cost x7.6, power x7.9
out.append("const CLICK_UPG=[")
for i,n in enumerate(UPG_NAMES):
    base = 25 * (7.0 ** i)
    add  = 1 * (7.3 ** i)
    out.append("{id:%s,n:%s,ico:%s,era:%d,base:%s,g:1.18,add:%s}," % (
        j("cu%d"%(i+1)), j(n), j(UPG_ICONS[i]), i//10, repr(base), repr(add)))
out.append("];")

# ---- CASES: cost x9.0 to track bot income; odds shift to rare with index
out.append("const CASES=[")
for i,n in enumerate(CASE_NAMES):
    cost = 500 * (7.8 ** i)
    w = [50*(0.965**i), 28*(0.985**i), 15*(1.0**i), 6*(1.012**i),
         2.5*(1.022**i), 1.0*(1.032**i), 0.12*(1.042**i)]
    tot = sum(w)
    odds = [round(x/tot, 6) for x in w]
    out.append("{id:%s,n:%s,art:%s,era:%d,cost:%s,base:%s,odds:%s}," % (
        j("c%d"%(i+1)), j(n+" Case"), j(CASE_ART[i//10]), i//10,
        repr(cost), repr(cost/6.0), j(odds)))
out.append("];")

out.append("const ERAS=%s;" % j(ERAS))
out.append("const CASE_ERAS=%s;" % j(CASE_ERAS))

open("content.js","w").write("\n".join(out))

# ---------------------------------------------------------------- sanity report
print("bots   : %d   first cost %.0f  last cost %.3g   last cps %.3g"
      % (len(BOT_NAMES), 15, 15*7.8**99, 0.15*8.3**99))
print("upgrades: %d  first cost %.0f  last cost %.3g   last add %.3g"
      % (len(UPG_NAMES), 25, 25*7.0**99, 7.3**99))
print("cases  : %d   first cost %.0f  last cost %.3g"
      % (len(CASE_NAMES), 500, 500*7.8**99))
w0=[50,28,15,6,2.5,1,0.12]; t0=sum(w0)
i=99
w9=[50*(0.965**i),28*(0.985**i),15,6*(1.012**i),2.5*(1.022**i),1.0*(1.032**i),0.12*(1.042**i)]
t9=sum(w9)
print("rare odds: first case %.3f%%   last case %.2f%%" % (100*w0[6]/t0, 100*w9[6]/t9))
print("names unique:", len(set(BOT_NAMES))==100, len(set(UPG_NAMES))==100, len(set(CASE_NAMES))==100)
