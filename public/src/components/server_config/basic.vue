<template>
    <collapsible :title="$t('title')">
        <field type="number" :label="$t('udp_label')" v-model="udpPort"></field>
        <field type="number" :label="$t('tcp_label')" v-model="tcpPort"></field>
        <field type="number" :label="$t('maxconnections_label')" v-model="maxConnections"></field>
        <checkbox :label="$t('registertolobby_label')" v-model="registerToLobby"></checkbox>
        <checkbox :label="$t('landiscovery_label')" v-model="lanDiscovery"></checkbox>
    </collapsible>
</template>

<script>
import collapsible from "../collapsible.vue";
import field from "../field.vue";
import checkbox from "../checkbox.vue";

export default {
    components: {collapsible, field, checkbox},
    data() {
        return {
            udpPort: 9600,
            tcpPort: 9600,
            maxConnections: 10,
            registerToLobby: true,
            lanDiscovery: false
        };
    },
    methods: {
        setData(data) {
            this.udpPort = data.udpPort;
            this.tcpPort = data.tcpPort;
            this.maxConnections = data.maxConnections;
            this.registerToLobby = data.registerToLobby;
            this.lanDiscovery = data.lanDiscovery;
        },
        getData() {
            return {
                udpPort: parseInt(this.udpPort),
                tcpPort: parseInt(this.tcpPort),
                maxConnections: parseInt(this.maxConnections),
                registerToLobby: this.registerToLobby ? 1 : 0,
                lanDiscovery: this.lanDiscovery ? 1 : 0
            };
        }
    }
}
</script>

<i18n>
{
    "en": {
        "title": "Basic configuration",
        "udp_label": "UDP port",
        "tcp_label": "TCP port",
        "maxconnections_label": "Max. connections",
        "registertolobby_label": "Register to lobby",
        "landiscovery_label": "LAN Discovery"
    }
}
</i18n>
