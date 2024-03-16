# General

* MC is a spacer working on an inter-system freighter.

* MC gets standard pay based on their position and rating.
* MC gets bonus pay based on ship profits based on their position and rating.
* MC gets a mass allotment based on their position and rating.

* MC can buy more mass allotment on the freighter.
* MC can buy and sell goods at every port.
* MC can look up information in a galactic world book about each port, planet, and station.
* MC can look up estimated price and demand for goods at any port.
  * Accuracy of estimated price and demand decreases significantly the further MC is from that port.

* MC does not determine where the freighter goes.

# Implementation Order
1. Location
   1. Planet
   2. Moon
   3. Space Habitat
   4. Mine
   5. Dock
   6. Industrial Facilities
      1. Shipyard
      2. Research lab
      3. Processing facility
      4. Fabrication facility
2. Location: Coordinates/Position
   1. Research
3. Freighter
4. Freighter: Coordinates/Position
   1. Research
5. Freighter: Velocity
6. Freighter: Max Speed
7. Freighter: Max Acceleration
8. Random Captain
   1. Captain randomly picks one of the 3 to 5 nearest locations to travel to
9. Location:

10. MC
11. MC: Rating
12. MC: Base Pay
13. MC: Share Pay
14. Freighter: Available Mass
15. MC: Mass Allotment


9. Location: Types
   1. planet w/space dock
   2. Space habitat
       1. https://en.wikipedia.org/wiki/Space_habitat
