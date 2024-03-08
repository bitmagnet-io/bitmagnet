export type Status = {
    status: string;
    timestamp: string;
    component: StatusComponent;
};

export type StatusComponent = {
    name: string;
    version: string;
};