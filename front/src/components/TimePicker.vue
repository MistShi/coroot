<template>
    <v-menu close-on-content-click offset-y>
        <template #activator="{ on, attrs }">
            <v-btn v-on="on" plain outlined height="40" class="px-2">
                <v-icon>mdi-clock-outline</v-icon>
                <span v-if="!small" class="ml-2">{{selected.text}}</span>
                <v-icon v-if="!small" small class="ml-2">
                    mdi-chevron-{{attrs['aria-expanded'] === 'true' ? 'up' : 'down'}}
                </v-icon>
            </v-btn>
        </template>
        <v-list dense dark>
            <template v-if="small && selected.custom">
                <v-list-item class="d-block text-center">
                    <div>{{selected.from}}</div>
                    <div>to</div>
                    <div>{{selected.to}}</div>
                </v-list-item>
                <v-divider />
            </template>
            <v-list-item v-for="i in intervals" :key="i.text" @click="changeInterval(i.query)" :input-value="selected.text === i.text">
                {{i.text}}
            </v-list-item>
        </v-list>
    </v-menu>
</template>

<script>
export default {
    props: {
        small: Boolean,
    },

    computed: {
        intervals() {
            return [
                {text: 'last hour', query: {from: 'now-1h', to: 'now'}},
                {text: 'last 3 hours', query: {from: 'now-3h', to: 'now'}},
                {text: 'last 12 hours', query: {from: 'now-12h', to: 'now'}},
                {text: 'last day', query: {from: 'now-1d', to: 'now'}},
                {text: 'last week', query: {from: 'now-7d', to: 'now'}},
            ];
        },
        selected() {
            const from = this.$route.query.from || 'now-1h';
            const to = this.$route.query.to || 'now';
            const v = this.intervals.find((i) => i.query.from === from && i.query.to === to);
            if (v) {
                return {text: v.text, custom: false};
            }
            const iFrom = parseInt(from);
            const iTo = parseInt(to);
            const format = (t) => this.$format.date(t, '{YYYY}-{MM}-{DD} {HH}:{mm}');
            const f = isNaN(iFrom) ? from : format(iFrom);
            const t = isNaN(iTo) ? to : format(iTo);
            return {from: f, to: t, text: f+' to '+t, custom: true};
        },
    },

    methods: {
        changeInterval(query) {
            this.$router.push({query: {...this.$route.query, ...query}}).catch(err => err);
        },
    },
}
</script>

<style scoped>
</style>
