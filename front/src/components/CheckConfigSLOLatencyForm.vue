<template>
    <div>
        Histogram query:
        <MetricSelector v-model="config.histogram_query" :rules="[$validators.notEmpty]" wrap="sum by(le)( rate( <input> [..]) )" class="mb-3"/>

        Objective:
        <div>
            <v-text-field outlined dense v-model="config.objective_percentage" :rules="[$validators.notEmpty]" hide-details class="input text">
                <template #append><span class="grey--text">%</span></template>
            </v-text-field>
            of requests should be served faster than
            <v-select v-model="config.objective_bucket" :items="buckets" :rules="[$validators.notEmpty]" outlined dense hide-details :menu-props="{offsetY: true}" class="input select" />
        </div>
    </div>
</template>

<script>
import MetricSelector from "@/components/MetricSelector";

const buckets = [
    {value: '0.005', text: '5ms'},
    {value: '0.01', text: '10ms'},
    {value: '0.025', text: '25ms'},
    {value: '0.05', text: '50ms'},
    {value: '0.1', text: '100ms'},
    {value: '0.25', text: '250ms'},
    {value: '0.5', text: '500ms'},
    {value: '1', text: '1s'},
    {value: '2.5', text: '2.5s'},
    {value: '5', text: '5s'},
    {value: '10', text: '10s'},
];
export default {
    components: {MetricSelector},
    props: {
        form: Object,
    },
    computed: {
        config() {
            return this.form.configs[0];
        },
        buckets() {
            return buckets;
        },
    },
}
</script>

<style scoped>
.input {
    display: inline-flex;
}
.input >>> .v-input__slot {
    min-height: initial !important;
    height: 1.5rem !important;
    padding: 0 8px !important;
}
.input.text >>> .v-input__append-inner {
    margin-top: 4px !important;
}
.input.select >>> .v-input__append-inner {
    margin-top: 0 !important;
}
.input >>> .v-select__selection--comma {
    margin: 0 !important;
}
* >>> .v-list-item {
    min-height: 32px !important;
    padding: 0 8px !important;
}
.input.text {
    max-width: 7ch;
}
.input.select {
    max-width: 11ch;
}
</style>