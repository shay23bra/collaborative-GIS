const map = L.map('map').setView([32.0853, 34.7818], 13);
let socket = connectWebSocket();

const openStreetMapLayer = L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
});
const existingLayerGroup = L.layerGroup().addTo(map);

// const satelliteLayer = L.tileLayer('https://{s}.govmap.gov.il/tiles/1.0.0/satellite/{z}/{x}/{y}.jpg', {
//     attribution: '&copy; govmap.gov.il',
//     maxZoom: 18,
// });

const satelliteLayer = L.tileLayer(
    'https://server.arcgisonline.com/ArcGIS/rest/services/World_Imagery/MapServer/tile/{z}/{y}/{x}',
    {
        attribution: 'Tiles &copy; Esri &mdash; Source: Esri, i-cubed, USDA, USGS, AEX, GeoEye, Getmapping, Aerogrid, IGN, IGP, UPR-EGP, and the GIS User Community',
        maxZoom: 18
    }
);


openStreetMapLayer.addTo(map);

L.control.layers({
    "OpenStreetMap": openStreetMapLayer,
    "Satellite": satelliteLayer
}).addTo(map);

const drawnItems = new L.FeatureGroup();
map.addLayer(drawnItems);

const drawControl = new L.Control.Draw({
    edit: {
        featureGroup: drawnItems
    },
    draw: {
        polygon: true,
        polyline: false,
        rectangle: false,
        circle: false,
        marker: false,
        circlemarker: false
    }
});
map.addControl(drawControl);

socket.onopen = () => {
    console.log("WebSocket connection established");
};

socket.onerror = (error) => {
    console.error("WebSocket error:", error);
};

function sendPolygonData(layer) {
    const coordinates = layer.getLatLngs()[0].map(latlng => [latlng.lat, latlng.lng]);
    const username = window.localStorage.getItem("username");
    const data = JSON.stringify({
        type: 'polygon',
        username: username,
        coordinates: coordinates
    });

    if (socket.readyState === WebSocket.OPEN) {
        socket.send(data);
    } else {
        console.error("WebSocket is not open. Cannot send data.");
    }
}

map.on(L.Draw.Event.CREATED, (event) => {
    const layer = event.layer;
    drawnItems.addLayer(layer);
    sendPolygonData(layer);
});

// map.on('load', loadAreasWithinBounds);
// map.on('moveend', loadAreasWithinBounds);


function connectWebSocket() {
    const backendUrl = "http://localhost:8000";
    const wsProtocol = backendUrl.startsWith("https") ? "wss" : "ws";
    const socket = new WebSocket(`${wsProtocol}://${backendUrl.replace(/^https?:\/\//, '')}/ws`);


    socket.onopen = () => {
        console.log("WebSocket connection established");
    };

    socket.onclose = () => {
        console.log("WebSocket connection closed. Reconnecting...");
        setTimeout(connectWebSocket, 1000);
    };

    socket.onerror = (error) => {
        console.error("WebSocket error:", error);
    };

    socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        if (data.type === 'polygon' && data.username !== username) { // Show only other users' drawings
            const polygon = L.polygon(data.coordinates);
            polygon.bindPopup(`Drawn by ${data.username}`).openPopup();
            drawnItems.addLayer(polygon);
        }
    };

    return socket;
}

map.on(L.Draw.Event.CREATED, (event) => {
    const layer = event.layer;
    drawnItems.addLayer(layer);

    const area = L.GeometryUtil.geodesicArea(layer.getLatLngs()[0]);
    layer.bindPopup(`Area: ${(area / 1e6).toFixed(2)} km²`).openPopup();

    sendPolygonData(layer);
});

socket.onmessage = (event) => {
    const data = JSON.parse(event.data);
    if (data.type === 'user_count') {
        document.getElementById("active-users").innerText = `Active users: ${data.count}`;
    }
};

socket.onclose = (event) => {
    console.log("WebSocket closed:", event);
    alert("Connection lost. Trying to reconnect...");
    setTimeout(() => {
        socket = connectWebSocket(); 
    }, 5000);
};

async function loadAreasWithinBounds() {
    // Get current map bounds
    const bounds = map.getBounds();
    const minLat = bounds.getSouth();
    const minLng = bounds.getWest();
    const maxLat = bounds.getNorth();
    const maxLng = bounds.getEast();

    try {
        const response = await fetch(`localhost:8000/areas_in_bounds?minLat=${minLat}&minLng=${minLng}&maxLat=${maxLat}&maxLng=${maxLng}`);
        if (response.ok) {
            const areas = await response.json();

            // Clear any existing layers before adding new ones
            existingLayerGroup.clearLayers();

            // Loop through each area and add it to the map
            areas.forEach(area => {
                const username = area.username;
                const geoJsonData = area.polygon; // Assuming the polygon is in GeoJSON format

                // Create a Leaflet GeoJSON layer for each area
                const geoJsonLayer = L.geoJSON(geoJsonData, {
                    onEachFeature: function (feature, layer) {
                        layer.bindPopup(`Created by: ${username}`);
                    }
                });

                // Add each GeoJSON layer to the layer group
                geoJsonLayer.addTo(existingLayerGroup);
            });
        } else {
            console.error("Failed to load areas:", response.statusText);
        }
    } catch (error) {
        console.error("Error loading areas within bounds:", error);
    }
}